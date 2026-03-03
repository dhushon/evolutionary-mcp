import { useState, useCallback, useContext, useEffect } from 'react';
import ReactFlow, {
  Node,
  Edge,
  addEdge,
  useNodesState,
  useEdgesState,
  Controls,
  Background,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { LogContext } from '../context/LogContext';
import FilterNode from './nodes/FilterNode';
import TransformNode from './nodes/TransformNode';
import WorkflowForm from './WorkflowForm';
import { useWorkflow, usePutWorkflow } from '../hooks/useWorkflows';
import { WorkflowUpdatePayload } from '../types';

const nodeTypes = {
    filter: FilterNode,
    transform: TransformNode,
};

const initialNodes: Node[] = [
  { id: '1', type: 'input', data: { label: 'Input Node' }, position: { x: 250, y: 5 } },
];
const initialEdges: Edge[] = [];
let idCounter = 2;
const getNextId = () => `${idCounter++}`;

interface WorkflowEditorProps {
    workflowId: string | null;
    isCreating: boolean;
    onSave?: () => void;
}

export default function WorkflowEditor({ workflowId, isCreating, onSave }: WorkflowEditorProps) {
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
    const [showSettings, setShowSettings] = useState(isCreating);
    
    const logContext = useContext(LogContext);
    const { data: workflow, isLoading: workflowLoading } = useWorkflow(workflowId);
    const putWorkflowMutation = usePutWorkflow();

    // Update form when workflow data loads
    useEffect(() => {
        if (workflow) {
            logContext?.addLog(`Loaded workflow: ${workflow.name} (v${workflow.version})`);
            // Here you would typically also load nodes/edges from a field in the workflow 
            // if they were stored in the DB. For now, we keep the UI state.
        }
    }, [workflow, logContext]);

    const handleSave = async (payload: WorkflowUpdatePayload) => {
        try {
            logContext?.addLog(`Saving ${payload.save_as_new_version ? 'new version' : 'draft'}...`);
            await putWorkflowMutation.mutateAsync(payload);
            logContext?.addLog('Save successful!');
            setShowSettings(false);
            if (onSave) onSave();
        } catch (err) {
            logContext?.addLog(`Save failed: ${err}`);
        }
    };

    const onConnect = useCallback(
        (params: any) => {
            logContext?.addLog(`Connected node ${params.source} to ${params.target}`);
            setEdges((eds) => addEdge(params, eds));
        },
        [setEdges, logContext],
    );

    const onAdd = useCallback((type: string) => {
        const newNodeId = getNextId();
        const newNode = {
            id: newNodeId,
            type,
            data: { label: `${type} node` },
            position: {
                x: 100 + Math.random() * 400,
                y: 100 + Math.random() * 400,
            },
        };
        logContext?.addLog(`Added ${type} node ${newNodeId}`);
        setNodes((nds) => nds.concat(newNode));
    }, [setNodes, logContext]);

    const onRun = useCallback(() => {
        logContext?.addLog('Running workflow...');
        const inputNode = nodes.find(node => node.type === 'input');
        if (!inputNode) {
            logContext?.addLog('Error: No input node found.');
            return;
        }

        let currentNode = inputNode;
        const visited = new Set<string>();

        const traverse = (node: Node) => {
            if (visited.has(node.id)) return;
            
            logContext?.addLog(`Executing node ${node.id}: ${node.data.label}`);
            visited.add(node.id);

            const outgoingEdges = edges.filter(edge => edge.source === node.id);
            outgoingEdges.forEach(edge => {
                const nextNode = nodes.find(n => n.id === edge.target);
                if (nextNode) {
                    traverse(nextNode);
                }
            });
        };

        traverse(currentNode);
        logContext?.addLog('Workflow finished.');
    }, [nodes, edges, logContext]);

    if (workflowId && workflowLoading) {
        return (
            <div className="flex items-center justify-center h-full">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
        );
    }
    
    return (
        <div className="relative h-full w-full bg-bg-base rounded-xl overflow-hidden border border-border-base">
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                nodeTypes={nodeTypes}
                fitView
            >
                <div className="absolute top-4 right-4 z-10 flex gap-2">
                    <button 
                        onClick={() => setShowSettings(!showSettings)} 
                        className="bg-white dark:bg-gray-800 text-text-base border border-border-base px-4 py-2 rounded-lg shadow-sm hover:bg-bg-accent transition-all flex items-center space-x-2"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        <span>{showSettings ? 'Hide Settings' : 'Settings'}</span>
                    </button>
                    <button onClick={() => onAdd('filter')} className="bg-white dark:bg-gray-800 text-text-base border border-border-base px-4 py-2 rounded-lg shadow-sm hover:bg-bg-accent transition-all">
                        + Filter
                    </button>
                    <button onClick={() => onAdd('transform')} className="bg-white dark:bg-gray-800 text-text-base border border-border-base px-4 py-2 rounded-lg shadow-sm hover:bg-bg-accent transition-all">
                        + Transform
                    </button>
                    <button onClick={onRun} className="bg-primary text-white px-4 py-2 rounded-lg shadow-sm hover:bg-primary-hover transition-all flex items-center space-x-2">
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span>Run</span>
                    </button>
                </div>

                {showSettings && (
                    <div className="absolute top-16 right-4 z-20 w-96 max-h-[80%] overflow-y-auto animate-in fade-in slide-in-from-top-2 duration-200">
                        <WorkflowForm
                            initialData={workflow || {
                                name: '',
                                description: '',
                                element_type: 'workflow',
                                status: 'draft',
                                input_schema: {},
                                output_schema: {}
                            }}
                            onSave={handleSave}
                            onCancel={() => setShowSettings(false)}
                            isSubmitting={putWorkflowMutation.isPending}
                        />
                    </div>
                )}

                <Controls />
                <Background />
            </ReactFlow>
        </div>
    );
}
