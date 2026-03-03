import React, { useMemo, useState } from 'react';
import { useWorkflows, useTenant } from '../hooks/useWorkflows';
import StatusBadge from './StatusBadge';
import { Workflow } from '../types';
import { Link, useLocation, useNavigate } from 'react-router-dom';

interface SidebarProps {
    selectedId: string | null;
    onSelect: (id: string) => void;
    onCreate: () => void;
    children?: React.ReactNode;
}

interface WorkflowNode extends Workflow {
    children: WorkflowNode[];
}

interface TreeNodeProps {
    node: WorkflowNode;
    selectedId: string | null;
    onSelect: (id: string) => void;
    level: number;
}

const TreeNode: React.FC<TreeNodeProps> = ({ node, selectedId, onSelect, level }) => {
    const [isExpanded, setIsExpanded] = useState(true);
    const isSelected = selectedId === node.id;
    const hasChildren = node.children.length > 0;

    return (
        <div className="flex flex-col">
            <button
                onClick={() => {
                    onSelect(node.id);
                    if (hasChildren) setIsExpanded(!isExpanded);
                }}
                className={`w-full flex items-center space-x-2 px-2 py-2 text-sm rounded-md transition-all duration-200 group ${
                    isSelected 
                        ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 font-medium' 
                        : 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/50'
                }`}
                style={{ paddingLeft: `${level * 12 + 8}px` }}
            >
                {hasChildren ? (
                    <svg 
                        className={`w-4 h-4 transition-transform duration-200 ${isExpanded ? 'rotate-90 text-blue-500' : 'text-slate-400 group-hover:text-slate-500'}`} 
                        fill="none" 
                        viewBox="0 0 24 24" 
                        stroke="currentColor"
                    >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                ) : (
                    <div className="w-4 h-4" /> // Spacer
                )}
                
                <span className="truncate flex-1 text-left">{node.name}</span>
                
                {node.status === 'active' && (
                    <span className="w-1.5 h-1.5 rounded-full bg-green-500"></span>
                )}
            </button>
            {hasChildren && isExpanded && (
                <div className="flex flex-col">
                    {node.children.map(child => (
                        <TreeNode 
                            key={child.id} 
                            node={child} 
                            selectedId={selectedId} 
                            onSelect={onSelect} 
                            level={level + 1} 
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default function Sidebar({ selectedId, onSelect, onCreate, children }: SidebarProps) {
    const { data: workflows, isLoading: workflowsLoading } = useWorkflows();
    const { data: tenant } = useTenant();
    const location = useLocation();
    const navigate = useNavigate();

    const workflowTree = useMemo(() => {
        if (!workflows) return [];
        const map = new Map<string, WorkflowNode>();
        const tree: WorkflowNode[] = [];

        workflows.forEach(wf => map.set(wf.id, { ...wf, children: [] }));
        workflows.forEach(wf => {
            const node = map.get(wf.id);
            if (!node) return;
            if (wf.parent_id && map.has(wf.parent_id)) {
                map.get(wf.parent_id)!.children.push(node);
            } else {
                tree.push(node);
            }
        });
        return tree;
    }, [workflows]);

    return (
        <aside className="w-64 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col h-full shadow-sm z-10">
            {/* 1) Company Identification */}
            <div className="p-5 border-b border-slate-200 dark:border-slate-800 flex items-center space-x-3">
                {tenant?.logo_svg ? (
                    <div className="w-8 h-8 flex items-center justify-center" dangerouslySetInnerHTML={{ __html: tenant.logo_svg }} />
                ) : (
                    <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center shadow-sm text-white font-bold">
                        {tenant?.name?.charAt(0) || 'E'}
                    </div>
                )}
                <div>
                    <h1 className="text-sm font-bold text-slate-900 dark:text-white leading-tight">
                        {tenant?.brand_title || tenant?.name || 'Evolutionary'}
                    </h1>
                    <p className="text-xs text-slate-500 dark:text-slate-400">Memory MCP</p>
                </div>
            </div>
            
            {/* 2) Main Navigation */}
            <div className="px-2 py-4 border-b border-slate-200 dark:border-slate-800 space-y-1">
                <Link 
                    to="/grounding"
                    className={`flex items-center space-x-3 px-3 py-2 rounded-lg text-sm transition-all ${
                        location.pathname === '/grounding' 
                            ? 'bg-blue-600 text-white shadow-md' 
                            : 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/50'
                    }`}
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
                    </svg>
                    <span className="font-medium">Grounding Rules</span>
                </Link>
                <Link 
                    to="/memories"
                    className={`flex items-center space-x-3 px-3 py-2 rounded-lg text-sm transition-all ${
                        location.pathname === '/memories' 
                            ? 'bg-blue-600 text-white shadow-md' 
                            : 'text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/50'
                    }`}
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                    </svg>
                    <span className="font-medium">Memory Inspector</span>
                </Link>
            </div>

            {children}

            {/* 3) Workflows Section */}
            <div className="flex-1 overflow-y-auto py-4">
                <div className="px-4 mb-2 flex items-center justify-between group">
                    <h2 className="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
                        Workflows
                    </h2>
                    <Link 
                        to="/workflows"
                        onClick={onCreate}
                        className="p-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-400 hover:text-blue-600 transition-colors"
                        title="Create Workflow"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                    </Link>
                </div>

                <nav className="space-y-0.5 px-2">
                    {workflowTree.map((node) => (
                        <TreeNode 
                            key={node.id} 
                            node={node} 
                            selectedId={selectedId} 
                            onSelect={(id) => {
                                onSelect(id);
                                // Ensure we are on the workflows page
                                if(location.pathname !== '/workflows' && location.pathname !== '/') {
                                    navigate('/workflows');
                                }
                            }} 
                            level={0} 
                        />
                    ))}
                    
                    {(!workflowsLoading && workflowTree.length === 0) && (
                        <div className="px-8 py-4 text-xs text-slate-400 italic">
                            No workflows found.
                        </div>
                    )}
                </nav>
            </div>

            {/* 4) Status Flags */}
            <StatusBadge />
        </aside>
    );
}
