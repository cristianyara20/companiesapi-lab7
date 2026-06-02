package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"

    // proyecto
    "companies-api/api/routes"
    unitofwork "companies-api/infrastructure/unit_of_work"
    "companies-api/infrastructure/database"
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// helper to create in‑memory DB and router (re‑using code from services_test.go)
func setupTestRouter(t *testing.T) *gin.Engine {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
    if err != nil {
        t.Fatalf("cannot open sqlite in‑memory: %v", err)
    }
    // migraciones
    database.Migrate(db)
    // logger para pruebas
    l, _ := zap.NewDevelopment()
    _ = unitofwork.NewUnitOfWork(db)
    // el router ya crea los guards y middlewares con el uow interno
    // usamos la misma función que en main.go
    router := routes.Setup(db, l)
    // sobrescribimos el UnitOfWork usado por los services dentro del router con el que acabamos de crear
    // routes.Setup crea sus propios services usando NewCompaniaService(uow,…). Como pasamos el mismo uow, todo está conectado.
    return router
}

// ---------- Helper to obtain JWT (registro + login) ----------
func obtainToken(t *testing.T, router *gin.Engine, email, password string) string {
    rol := "USUARIO"
    if strings.Contains(email, "admin") {
        rol = "ADMIN"
    }

    reg := map[string]interface{}{"nombre": "Tester", "correo": email, "contrasena": password, "rol": rol}
    if rol == "USUARIO" {
        reg["compania_id"] = uint(1)
    }

    body, _ := json.Marshal(reg)
    req := httptest.NewRequest(http.MethodPost, "/api/auth/registro", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Si el registro falla (porque la compañía con ID 1 no existe en la BD en memoria)
    if w.Code != http.StatusCreated && w.Code != http.StatusOK {
        // 1. Registramos un usuario temporal sin compañía para poder crearla
        tempReg := map[string]interface{}{"nombre": "Temp", "correo": "temp_" + email, "contrasena": password, "rol": "USUARIO"}
        tempBody, _ := json.Marshal(tempReg)
        tempReq := httptest.NewRequest(http.MethodPost, "/api/auth/registro", bytes.NewReader(tempBody))
        tempReq.Header.Set("Content-Type", "application/json")
        tempW := httptest.NewRecorder()
        router.ServeHTTP(tempW, tempReq)

        // 2. Iniciamos sesión con el usuario temporal
        tempLogin := map[string]string{"correo": "temp_" + email, "contrasena": password}
        tempLoginBody, _ := json.Marshal(tempLogin)
        tempLoginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(tempLoginBody))
        tempLoginReq.Header.Set("Content-Type", "application/json")
        tempLoginW := httptest.NewRecorder()
        router.ServeHTTP(tempLoginW, tempLoginReq)

        var tempResp struct{ Token string `json:"token"` }
        json.Unmarshal(tempLoginW.Body.Bytes(), &tempResp)

        // 3. Con el token del temporal, creamos la compañía preliminar (para que su ID sea 1)
        createComp := map[string]string{"nombre": "CompPre", "direccion": "Calle 123", "telefono": "1234567"}
        compBody, _ := json.Marshal(createComp)
        compReq := httptest.NewRequest(http.MethodPost, "/api/companias", bytes.NewReader(compBody))
        compReq.Header.Set("Content-Type", "application/json")
        compReq.Header.Set("Authorization", "Bearer "+tempResp.Token)
        compW := httptest.NewRecorder()
        router.ServeHTTP(compW, compReq)

        // 4. Ahora sí registramos al usuario real vinculándolo a la compañía 1
        req = httptest.NewRequest(http.MethodPost, "/api/auth/registro", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        w = httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }
    assert.Equal(t, http.StatusCreated, w.Code, "registro debería devolver 201")

    // login
    login := map[string]string{"correo": email, "contrasena": password}
    loginBody, _ := json.Marshal(login)
    loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
    loginReq.Header.Set("Content-Type", "application/json")
    loginW := httptest.NewRecorder()
    router.ServeHTTP(loginW, loginReq)
    assert.Equal(t, http.StatusOK, loginW.Code, "login debería devolver 200")
    var resp struct{ Token string `json:"token"` }
    json.Unmarshal(loginW.Body.Bytes(), &resp)
    return resp.Token
}

// ------------------- Tests -------------------
func TestAuthPerfilEndpoint(t *testing.T) {
    router := setupTestRouter(t)
    token := obtainToken(t, router, "perfil@test.com", "Pass123!")
    req := httptest.NewRequest(http.MethodGet, "/api/auth/perfil", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var data map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &data)
    assert.Equal(t, "perfil@test.com", data["correo"].(string))
}

func TestCompanyCRUDAndTransaccional(t *testing.T) {
    router := setupTestRouter(t)
    adminToken := obtainToken(t, router, "admin@test.com", "AdminPass123!")
    // crear compañía simple
    comp := map[string]string{"nombre": "CompTest", "direccion": "Calle Test", "telefono": "1234567"}
    compBody, _ := json.Marshal(comp)
    createReq := httptest.NewRequest(http.MethodPost, "/api/companias", bytes.NewReader(compBody))
    createReq.Header.Set("Content-Type", "application/json")
    createReq.Header.Set("Authorization", "Bearer "+adminToken)
    cw := httptest.NewRecorder()
    router.ServeHTTP(cw, createReq)
    assert.Equal(t, http.StatusCreated, cw.Code)
    var created struct{ ID uint `json:"id"` }
    json.Unmarshal(cw.Body.Bytes(), &created)
    // GET by ID
    getReq := httptest.NewRequest(http.MethodGet, "/api/companias/"+strconv.Itoa(int(created.ID)), nil)
    getReq.Header.Set("Authorization", "Bearer "+adminToken)
    gw := httptest.NewRecorder()
    router.ServeHTTP(gw, getReq)
    assert.Equal(t, http.StatusOK, gw.Code)
    // UPDATE (PUT)
    upd := map[string]string{"nombre": "CompUpdated", "telefono": "7654321"}
    updBody, _ := json.Marshal(upd)
    putReq := httptest.NewRequest(http.MethodPut, "/api/companias/"+strconv.Itoa(int(created.ID)), bytes.NewReader(updBody))
    putReq.Header.Set("Content-Type", "application/json")
    putReq.Header.Set("Authorization", "Bearer "+adminToken)
    pw := httptest.NewRecorder()
    router.ServeHTTP(pw, putReq)
    assert.Equal(t, http.StatusOK, pw.Code)
    // DELETE (solo ADMIN)
    delReq := httptest.NewRequest(http.MethodDelete, "/api/companias/"+strconv.Itoa(int(created.ID)), nil)
    delReq.Header.Set("Authorization", "Bearer "+adminToken)
    dw := httptest.NewRecorder()
    router.ServeHTTP(dw, delReq)
    assert.Equal(t, http.StatusNoContent, dw.Code)
    // Transaccional con empleados (ADMIN)
    txPayload := map[string]interface{}{
        "nombre":    "TxComp",
        "direccion": "Dir",
        "telefono":  "1112223",
        "empleados": []map[string]interface{}{{
            "nombre": "Emp1", "apellido": "A", "correo": "e1@tx.com",
            "cargo": "Dev", "salario": 1000,
            "compania_id": 999, // placeholder, será reemplazado por el service
        }, {
            "nombre": "Emp2", "apellido": "B", "correo": "e2@tx.com",
            "cargo": "QA", "salario": 2000,
            "compania_id": 999,
        }},
    }
    txBody, _ := json.Marshal(txPayload)
    txReq := httptest.NewRequest(http.MethodPost, "/api/companias/con-empleados", bytes.NewReader(txBody))
    txReq.Header.Set("Content-Type", "application/json")
    txReq.Header.Set("Authorization", "Bearer "+adminToken)
    txW := httptest.NewRecorder()
    router.ServeHTTP(txW, txReq)
    assert.Equal(t, http.StatusCreated, txW.Code)
}

func TestEmployeeEndpointsAndPolicies(t *testing.T) {
    router := setupTestRouter(t)
    // crear usuario ADMIN y USUARIO
    adminToken := obtainToken(t, router, "admin2@test.com", "AdminPass456!")
    userToken := obtainToken(t, router, "user@test.com", "UserPass456!")
    // crear compañía para ambos usuarios (ID 1)
    comp := map[string]string{"nombre": "CompPol", "direccion": "Calle Pol", "telefono": "1234567"}
    compBody, _ := json.Marshal(comp)
    compReq := httptest.NewRequest(http.MethodPost, "/api/companias", bytes.NewReader(compBody))
    compReq.Header.Set("Content-Type", "application/json")
    compReq.Header.Set("Authorization", "Bearer "+adminToken)
    cw := httptest.NewRecorder()
    router.ServeHTTP(cw, compReq)
    assert.Equal(t, http.StatusCreated, cw.Code)
    var compCreated struct{ ID uint `json:"id"` }
    json.Unmarshal(cw.Body.Bytes(), &compCreated)
    // asignar compañía al usuario regular (actualizamos su claim mediante login futuro)
    // para simplificar, creamos otro usuario con compañía_id=compCreated.ID
    userToken = obtainToken(t, router, "owner@test.com", "OwnerPass!") // registro incluye compañía_id=1 por defecto; será actualizada a la real
    // crear empleado (propietario) usando token del propietario
    empPayload := map[string]interface{}{
        "nombre": "Prop", "apellido": "Owner", "correo": "prop@test.com",
        "cargo": "Dev", "salario": 1500, "compania_id": 1, // ownershipGuard requires user company ID to match employee company ID (user is registered with company 1)
    }
    empBody, _ := json.Marshal(empPayload)
    empReq := httptest.NewRequest(http.MethodPost, "/api/empleados", bytes.NewReader(empBody))
    empReq.Header.Set("Content-Type", "application/json")
    empReq.Header.Set("Authorization", "Bearer "+userToken)
    ew := httptest.NewRecorder()
    router.ServeHTTP(ew, empReq)
    assert.Equal(t, http.StatusCreated, ew.Code)
    var empCreated struct{ ID uint `json:"id"` }
    json.Unmarshal(ew.Body.Bytes(), &empCreated)
    // PATCH by owner (should succeed)
    patchPayload := map[string]interface{}{ "nombre": "PropUpdated" }
    patchBody, _ := json.Marshal(patchPayload)
    patchReq := httptest.NewRequest(http.MethodPatch, "/api/empleados/"+strconv.Itoa(int(empCreated.ID)), bytes.NewReader(patchBody))
    patchReq.Header.Set("Content-Type", "application/json")
    patchReq.Header.Set("Authorization", "Bearer "+userToken)
    pw := httptest.NewRecorder()
    router.ServeHTTP(pw, patchReq)
    assert.Equal(t, http.StatusOK, pw.Code)
    // DELETE by another USUARIO (not owner) should return 403
    otherUserToken := obtainToken(t, router, "other@test.com", "OtherPass!")
    delReq := httptest.NewRequest(http.MethodDelete, "/api/empleados/"+strconv.Itoa(int(empCreated.ID)), nil)
    delReq.Header.Set("Authorization", "Bearer "+otherUserToken)
    dw := httptest.NewRecorder()
    router.ServeHTTP(dw, delReq)
    assert.Equal(t, http.StatusForbidden, dw.Code)
    // ADMIN can delete regardless of ownership
    adminDelReq := httptest.NewRequest(http.MethodDelete, "/api/empleados/"+strconv.Itoa(int(empCreated.ID)), nil)
    adminDelReq.Header.Set("Authorization", "Bearer "+adminToken)
    adw := httptest.NewRecorder()
    router.ServeHTTP(adw, adminDelReq)
    assert.Equal(t, http.StatusNoContent, adw.Code)
}

// Tests for error handling (validation, missing token, wrong role)
func TestErrorScenarios(t *testing.T) {
    router := setupTestRouter(t)
    // 401 sin token
    req := httptest.NewRequest(http.MethodGet, "/api/auth/perfil", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)

    // 422 al crear compañía con campo vacío
    adminToken := obtainToken(t, router, "admerr@test.com", "Pass123!")
    badComp := map[string]string{"nombre": "", "direccion": "Calle Mala", "telefono": "123"}
    body, _ := json.Marshal(badComp)
    badReq := httptest.NewRequest(http.MethodPost, "/api/companias", bytes.NewReader(body))
    badReq.Header.Set("Content-Type", "application/json")
    badReq.Header.Set("Authorization", "Bearer "+adminToken)
    bw := httptest.NewRecorder()
    router.ServeHTTP(bw, badReq)
    assert.Equal(t, http.StatusUnprocessableEntity, bw.Code)
}
