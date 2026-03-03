import React from 'react';
import LogsPanel from '../components/LogsPanel';
import WorkflowEditor from '../components/WorkflowEditor';
import { LogProvider } from '../context/LogContext';
import { useSelection } from '../context/SelectionContext';

const WorkflowExecution: React.FC = () => {
  const { selectedId, isCreating } = useSelection();

  return (
    <LogProvider>
      <div className="grid grid-rows-[auto_1fr] h-full gap-4 p-4 md:p-6 overflow-hidden">
        <div className="flex-none flex items-center justify-between">
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
        
        <div className="grid grid-rows-[1fr_auto] gap-4 min-h-0">
          <div className="min-h-0">
            <WorkflowEditor
              workflowId={selectedId}
              isCreating={isCreating}
            />
          </div>
          <div className="h-48 shrink-0">
            <LogsPanel />
          </div>
        </div>
      </div>
    </LogProvider>
  );
};

export default WorkflowExecution;
