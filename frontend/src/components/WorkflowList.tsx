import { useQuery } from '@tanstack/react-query';
import { getWorkflows } from '../api';

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
                    <div key={wf.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-5 shadow-sm hover:shadow-md transition-shadow">
                        <div className="flex justify-between items-start mb-3">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{wf.name}</h2>
                            <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                                wf.status === 'active' 
                                    ? 'bg-green-100 dark:bg-green-900/50 text-green-800 dark:text-green-300' 
                                    : 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
                            }`}>
                                {wf.status}
                            </span>
                        </div>
                        <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">{wf.description}</p>
                        <div className="flex justify-between items-center text-xs text-gray-500 dark:text-gray-400 border-t border-gray-200 dark:border-gray-700 pt-3">
                            <span className="font-mono bg-gray-100 dark:bg-gray-700 px-1.5 py-0.5 rounded">v{wf.version}</span>
                            <span>{new Date(wf.updated_at).toLocaleDateString()}</span>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}