import React, { useMemo, useState } from 'react';
import { useWorkflows, useTenant } from '../hooks/useWorkflows';
import StatusBadge from './StatusBadge';
import { Workflow } from '../types';
import { Link, useLocation } from 'react-router-dom';

interface SidebarProps {
    selectedId: string | null;
    onSelect: (id: string) => void;
    onCreate: () => void;
    children?: React.ReactNode;
}

interface TreeNodeProps {
    node: Workflow & { children: any[] };
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
                        ? 'bg-primary/10 text-primary dark:text-blue-400 font-medium' 
                        : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50'
                }`}
                style={{ paddingLeft: `${level * 12 + 8}px` }}
            >
                {hasChildren ? (
                    <svg 
                        className={`w-4 h-4 transition-transform duration-200 ${isExpanded ? 'rotate-90 text-primary' : 'text-gray-400 group-hover:text-gray-500'}`} 
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

    const workflowTree = useMemo(() => {
        if (!workflows) return [];
        const map = new Map<string, Workflow & { children: any[] }>();
        const tree: any[] = [];

        workflows.forEach(wf => map.set(wf.id, { ...wf, children: [] }));
        workflows.forEach(wf => {
            if (wf.parent_id && map.has(wf.parent_id)) {
                map.get(wf.parent_id)!.children.push(map.get(wf.id));
            } else {
                tree.push(map.get(wf.id));
            }
        });
        return tree;
    }, [workflows]);

    return (
        <aside className="w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col h-full shadow-sm z-10">
            {/* 1) Company Identification */}
            <div className="p-5 border-b border-gray-200 dark:border-gray-700 flex items-center space-x-3">
                {tenant?.logo_svg ? (
                    <div className="w-8 h-8 flex items-center justify-center" dangerouslySetInnerHTML={{ __html: tenant.logo_svg }} />
                ) : (
                    <div className="w-8 h-8 bg-gradient-to-br from-primary to-secondary rounded-lg flex items-center justify-center shadow-sm text-white font-bold">
                        {tenant?.name?.charAt(0) || 'E'}
                    </div>
                )}
                <div>
                    <h1 className="text-sm font-bold text-gray-900 dark:text-white leading-tight">
                        {tenant?.brand_title || tenant?.name || 'Evolutionary'}
                    </h1>
                    <p className="text-xs text-gray-500 dark:text-gray-400">Memory MCP</p>
                </div>
            </div>
            
            {/* 2) Main Navigation */}
            <div className="px-2 py-4 border-b border-gray-200 dark:border-gray-700">
                <Link 
                    to="/grounding"
                    className={`flex items-center space-x-3 px-3 py-2 rounded-lg text-sm transition-all ${
                        location.pathname === '/grounding' 
                            ? 'bg-primary text-white shadow-md' 
                            : 'text-text-base hover:bg-bg-accent'
                    }`}
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
                    </svg>
                    <span className="font-medium">Grounding Rules</span>
                </Link>
            </div>

            {children}

            {/* 3) Workflows Section */}
            <div className="flex-1 overflow-y-auto py-4">
                <div className="px-4 mb-2 flex items-center justify-between group">
                    <h2 className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Workflows
                    </h2>
                    <Link 
                        to="/workflows"
                        onClick={onCreate}
                        className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-400 hover:text-primary transition-colors"
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
                                    window.location.href = '/workflows';
                                }
                            }} 
                            level={0} 
                        />
                    ))}
                    
                    {(!workflowsLoading && workflowTree.length === 0) && (
                        <div className="px-8 py-4 text-xs text-gray-400 italic">
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
