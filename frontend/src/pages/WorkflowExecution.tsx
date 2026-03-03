import React from 'react';
import LogsPanel from '../components/LogsPanel';
import WorkflowEditor from '../components/WorkflowEditor';
import { LogProvider } from '../context/LogContext';
import { useSelection } from '../components/layout/DashboardLayout';

const WorkflowExecution: React.FC = () => {
  const { selectedId, isCreating } = useSelection();

  return (
    <LogProvider>
      <div className="flex flex-col h-full space-y-4 p-4 md:p-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-text-base">
            {isCreating ? 'Create Workflow Concept' : 'Workflow Designer'}
          </h1>
          {selectedId && (
            <div className="flex items-center space-x-2 text-sm text-text-muted">
              <span className="px-2 py-1 bg-bg-accent rounded border border-border-base font-mono">
                ID: {selectedId.substring(0, 8)}...
              </span>
            </div>
          )}
        </div>
        
        <div className="flex-1 flex flex-col min-h-0 space-y-4">
          <div className="flex-1 min-h-[400px]">
            <WorkflowEditor
              workflowId={selectedId}
              isCreating={isCreating}
            />
          </div>
          <div className="h-64 shrink-0">
            <LogsPanel />
          </div>
        </div>
      </div>
    </LogProvider>
  );
};

export default WorkflowExecution;
