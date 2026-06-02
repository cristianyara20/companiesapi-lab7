"use client";

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Plus, Edit2, Trash2, Users, Loader2, AlertCircle, Building2 } from 'lucide-react';

interface Empleado {
  id: number;
  nombre: string;
  apellido: string;
  correo: string;
  cargo: string;
  salario: number;
  compania_id: number;
}

interface Compania {
  id: number;
  nombre: string;
}

export default function EmpleadosPage() {
  const { user } = useAuth();
  const [empleados, setEmpleados] = useState<Empleado[]>([]);
  const [companias, setCompanias] = useState<Compania[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Modal state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({ id: 0, nombre: '', apellido: '', correo: '', cargo: '', salario: 0, compania_id: 0 });
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [empRes, compRes] = await Promise.all([
        api.get('/empleados?tamano=100&orden=id&dir=desc'),
        api.get('/companias')
      ]);
      setEmpleados(empRes.data.datos || []);
      setCompanias(compRes.data || []);
    } catch (err: any) {
      setError('Error al cargar la información');
    } finally {
      setLoading(false);
    }
  };

  const getCompaniaNombre = (id: number) => {
    return companias.find(c => c.id === id)?.nombre || 'Desconocida';
  };

  const openModal = (empleado?: Empleado) => {
    if (empleado) {
      setFormData(empleado);
      setIsEditing(true);
    } else {
      setFormData({ id: 0, nombre: '', apellido: '', correo: '', cargo: '', salario: 0, compania_id: user?.compania_id || 0 });
      setIsEditing(false);
    }
    setError('');
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setFormData({ id: 0, nombre: '', apellido: '', correo: '', cargo: '', salario: 0, compania_id: 0 });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Clean data payload, convert to float for salario, int for compania
      const payload = {
        ...formData,
        salario: parseFloat(formData.salario.toString()),
        compania_id: parseInt(formData.compania_id.toString(), 10)
      };

      if (isEditing) {
        // We use PUT for replacing, or PATCH for partial. Let's use PUT as per backend.
        await api.put(`/empleados/${formData.id}`, payload);
      } else {
        await api.post('/empleados', payload);
      }
      fetchData();
      closeModal();
    } catch (err: any) {
      if (err.response?.data?.errores) {
        const errorDetails = err.response.data.errores.map((e: any) => `${e.campo}: ${e.detalle}`).join(', ');
        setError(`Errores de validación: ${errorDetails}`);
      } else {
        setError(err.response?.data?.message || err.response?.data?.error || err.response?.data?.mensaje || 'Error al guardar el empleado');
      }
    }
  };

  const handlePartialUpdate = async (id: number) => {
    const newCargo = prompt('Introduce el nuevo cargo:');
    if (!newCargo) return;
    try {
      await api.patch(`/empleados/${id}`, { cargo: newCargo });
      fetchData();
    } catch (err: any) {
      alert(err.response?.data?.message || err.response?.data?.error || 'Error en actualización parcial (PATCH).');
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('¿Estás seguro de eliminar este empleado?')) {
      try {
        await api.delete(`/empleados/${id}`);
        fetchData();
      } catch (err: any) {
        alert(err.response?.data?.message || err.response?.data?.error || 'Error al eliminar el empleado (Posible falta de permisos).');
      }
    }
  };

  if (loading) return <div className="flex justify-center p-12"><Loader2 className="w-8 h-8 animate-spin text-purple-500" /></div>;

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex justify-between items-center bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-purple-100 flex items-center justify-center">
            <Users className="w-5 h-5 text-purple-600" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-slate-900">Gestión de Empleados</h1>
            <p className="text-sm text-slate-500">Directorio y roles del personal</p>
          </div>
        </div>
        <button 
          onClick={() => openModal()}
          className="bg-purple-600 hover:bg-purple-700 text-white px-5 py-2.5 rounded-xl text-sm font-medium transition-colors shadow-sm shadow-purple-500/20 flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          Nuevo Empleado
        </button>
      </div>

      {error && (
        <div className="bg-red-50 text-red-500 p-4 rounded-xl border border-red-100 flex items-center gap-2 text-sm">
          <AlertCircle className="w-5 h-5" />
          {error}
        </div>
      )}

      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200 text-sm font-semibold text-slate-600">
                <th className="px-6 py-4">ID</th>
                <th className="px-6 py-4">Nombre Completo</th>
                <th className="px-6 py-4">Cargo</th>
                <th className="px-6 py-4">Salario</th>
                <th className="px-6 py-4">Compañía</th>
                <th className="px-6 py-4 text-right">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {empleados.map((e) => (
                <tr key={e.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 text-sm font-medium text-slate-900">#{e.id}</td>
                  <td className="px-6 py-4">
                    <div className="flex flex-col">
                      <span className="text-sm font-medium text-slate-900">{e.nombre} {e.apellido}</span>
                      <span className="text-xs text-slate-500">{e.correo}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-sm text-slate-600">
                    <span className="bg-slate-100 text-slate-700 px-2 py-1 rounded-md text-xs font-medium">
                      {e.cargo}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm font-medium text-emerald-600">
                    ${e.salario.toLocaleString()}
                  </td>
                  <td className="px-6 py-4 text-sm text-slate-500 flex items-center gap-1.5 mt-2">
                    <Building2 className="w-3.5 h-3.5" />
                    {getCompaniaNombre(e.compania_id)}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center justify-end gap-1">
                      <button 
                        onClick={() => handlePartialUpdate(e.id)}
                        className="text-[10px] uppercase font-bold tracking-wider text-purple-600 hover:bg-purple-50 px-2 py-1.5 rounded-lg transition-colors border border-purple-100 mr-1"
                        title="Actualización Parcial (PATCH)"
                      >
                        Patch
                      </button>
                      <button 
                        onClick={() => openModal(e)}
                        className="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                        title="Editar completo (PUT)"
                      >
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button 
                        onClick={() => handleDelete(e.id)}
                        className="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        title="Eliminar"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
              {empleados.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-6 py-12 text-center text-slate-500 text-sm">
                    No hay empleados registrados
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-3xl w-full max-w-lg p-6 shadow-2xl border border-slate-200 animate-in zoom-in-95 duration-200">
            <h2 className="text-xl font-bold text-slate-900 mb-6">
              {isEditing ? 'Editar Empleado' : 'Nuevo Empleado'}
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Nombre</label>
                  <input
                    type="text"
                    value={formData.nombre}
                    onChange={(e) => setFormData({...formData, nombre: e.target.value})}
                    className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Apellido</label>
                  <input
                    type="text"
                    value={formData.apellido}
                    onChange={(e) => setFormData({...formData, apellido: e.target.value})}
                    className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all"
                    required
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Correo Electrónico</label>
                <input
                  type="email"
                  value={formData.correo}
                  onChange={(e) => setFormData({...formData, correo: e.target.value})}
                  className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Cargo</label>
                  <input
                    type="text"
                    value={formData.cargo}
                    onChange={(e) => setFormData({...formData, cargo: e.target.value})}
                    className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-1">Salario ($)</label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.salario}
                    onChange={(e) => setFormData({...formData, salario: Number(e.target.value)})}
                    className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Compañía (ID)</label>
                <select
                  value={formData.compania_id}
                  onChange={(e) => setFormData({...formData, compania_id: Number(e.target.value)})}
                  className="w-full border border-slate-200 rounded-xl px-4 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all bg-white"
                  required
                >
                  <option value="" disabled>Seleccione una compañía...</option>
                  {companias.map((c) => (
                    <option key={c.id} value={c.id}>{c.nombre} (ID: {c.id})</option>
                  ))}
                </select>
              </div>
              
              <div className="flex items-center justify-end gap-3 pt-4 mt-2 border-t border-slate-100">
                <button
                  type="button"
                  onClick={closeModal}
                  className="px-4 py-2 text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-xl transition-colors"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="bg-purple-600 hover:bg-purple-700 text-white px-5 py-2 rounded-xl text-sm font-medium transition-colors shadow-sm shadow-purple-500/20"
                >
                  Guardar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
