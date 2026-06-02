"use client";

import { useAuth } from '@/contexts/AuthContext';
import { Building2, Users, Activity, TrendingUp } from 'lucide-react';
import { useEffect, useState } from 'react';
import api from '@/lib/api';

export default function DashboardPage() {
  const { user } = useAuth();
  const [stats, setStats] = useState({ companias: 0, empleados: 0 });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [compRes, empRes] = await Promise.all([
          api.get('/companias'),
          api.get('/empleados')
        ]);
        setStats({
          companias: compRes.data.length || 0,
          empleados: empRes.data.total || 0
        });
      } catch (error) {
        console.error("Error fetching stats", error);
      }
    };
    fetchStats();
  }, []);

  const statCards = [
    { name: 'Total Compañías', value: stats.companias, icon: Building2, trend: '+12%', color: 'text-blue-600', bg: 'bg-blue-100' },
    { name: 'Total Empleados', value: stats.empleados, icon: Users, trend: '+4%', color: 'text-purple-600', bg: 'bg-purple-100' },
    { name: 'Actividad Reciente', value: 'Alta', icon: Activity, trend: 'Estable', color: 'text-emerald-600', bg: 'bg-emerald-100' },
  ];

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Hola, {user?.nombre} 👋</h1>
          <p className="text-slate-500 mt-1">Aquí está el resumen de tu cuenta hoy.</p>
        </div>
        <div className="flex items-center gap-2 bg-slate-100 px-4 py-2 rounded-xl border border-slate-200">
          <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse" />
          <span className="text-sm font-medium text-slate-700">Sistema Operativo</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {statCards.map((stat, i) => (
          <div key={i} className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow group cursor-default">
            <div className="flex justify-between items-start mb-4">
              <div className={`w-12 h-12 rounded-xl ${stat.bg} flex items-center justify-center group-hover:scale-110 transition-transform`}>
                <stat.icon className={`w-6 h-6 ${stat.color}`} />
              </div>
              <div className="flex items-center gap-1 text-emerald-600 text-sm font-medium bg-emerald-50 px-2 py-1 rounded-lg">
                <TrendingUp className="w-3.5 h-3.5" />
                {stat.trend}
              </div>
            </div>
            <h3 className="text-slate-500 text-sm font-medium mb-1">{stat.name}</h3>
            <p className="text-3xl font-bold text-slate-900">{stat.value}</p>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
        {/* Mock Chart/Recent Activity Area */}
        <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm h-80 flex flex-col justify-center items-center">
           <div className="w-16 h-16 rounded-full bg-slate-50 flex items-center justify-center mb-4">
             <Activity className="w-8 h-8 text-slate-300" />
           </div>
           <p className="text-slate-500 font-medium">Gráficos próximamente</p>
        </div>
        
        <div className="bg-white p-6 rounded-2xl border border-slate-200 shadow-sm h-80 flex flex-col justify-center items-center">
           <div className="w-16 h-16 rounded-full bg-slate-50 flex items-center justify-center mb-4">
             <Users className="w-8 h-8 text-slate-300" />
           </div>
           <p className="text-slate-500 font-medium">Actividad reciente próximamente</p>
        </div>
      </div>
    </div>
  );
}
