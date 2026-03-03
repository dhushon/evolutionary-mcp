import { useState } from 'react';

type Status = 'idle' | 'running' | 'success' | 'error';

const useWorkflowRunner = () => {
  const [status, setStatus] = useState<Status>('idle');
  const [logs, setLogs] = useState<string[]>([]);

  const runWorkflow = (_workflow: any) => {
    setStatus('running');
    setLogs(['Starting workflow...']);

    // Simulate workflow execution
    setTimeout(() => {
      setLogs(prev => [...prev, 'Step 1: Data extraction complete.']);
      setTimeout(() => {
        setLogs(prev => [...prev, 'Step 2: Data transformation complete.']);
        setTimeout(() => {
          setStatus('success');
          setLogs(prev => [...prev, 'Workflow finished successfully.']);
        }, 1000);
      }, 1000);
    }, 1000);
  };

  return { status, logs, runWorkflow };
};

export default useWorkflowRunner;
