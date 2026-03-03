import { useQuery } from '@tanstack/react-query';
import { getWorkflows } from '../api';
import { DataCard } from './layout/DataCard';

export default function WorkflowList() {
    const { data: workflows, isLoading, isError, error } = useQuery({
        queryKey: ['workflows'],
        queryFn: getWorkflows,
    });

    if (isLoading) return <div className="p-8 text-center text-gray-500 dark:text-gray-400">Loading workflows...</div>;
    if (isError) return <div className="p-8 text-center text-red-500">Error: {error.message}</div>;

    return (
        <div className="container mx-auto p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Workflows</h1>
                <button className="px-4 py-2 bg-primary text-white rounded-md hover:opacity-90 transition-opacity">
                    New Workflow
                </button>
            </div>
            
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {workflows?.map((wf) => (
                    <DataCard
                        key={wf.id}
                        title={wf.name}
                        subtitle={`v${wf.version} • ${wf.status}`}
                        content={wf.description}
                        footer={
                            <div className="flex justify-between items-center text-xs text-slate-500">
                                <span>Updated {new Date(wf.updated_at).toLocaleDateString()}</span>
                                <button className="text-blue-600 hover:text-blue-700 font-medium">View Details</button>
                            </div>
                        }
                    />
                ))}
            </div>
        </div>
    );
}