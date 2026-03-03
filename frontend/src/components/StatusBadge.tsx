import { useQuery } from '@tanstack/react-query';
import { getHealth } from '../api';
import { useEffect, useState } from 'react';

export default function StatusBadge() {
    const [lastChecked, setLastChecked] = useState<Date>(new Date());
    
    const { data: health, isError } = useQuery({
        queryKey: ['health'],
        queryFn: getHealth,
        refetchInterval: 30000, // Poll every 30s
    });

    useEffect(() => {
        if (health || isError) {
            setLastChecked(new Date());
        }
    }, [health, isError]);

    const isHealthy = health?.status === 'ok' && !isError;

    return (
        <div className="px-4 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
            <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                    <span className="relative flex h-2.5 w-2.5">
                        {isHealthy && (
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                        )}
                        <span className={`relative inline-flex rounded-full h-2.5 w-2.5 ${isHealthy ? 'bg-green-500' : 'bg-red-500'}`}></span>
                    </span>
                    <span className="text-xs font-medium text-gray-600 dark:text-gray-300">
                        {isHealthy ? 'System Online' : 'System Offline'}
                    </span>
                </div>
                <span className="text-[10px] text-gray-400">{lastChecked.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
            </div>
        </div>
    );
}