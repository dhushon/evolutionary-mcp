import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Shell } from './components/layout/Shell';
import WorkflowExecution from './pages/WorkflowExecution';
import GroundingManager from './pages/GroundingManager';
import MemoryInspector from './pages/MemoryInspector';
import { DataProvider } from './context/DataContext';

const App: React.FC = () => {
  return (
    <DataProvider>
      <Router>
        <Shell>
          <Routes>
            <Route path="/grounding" element={<GroundingManager />} />
            <Route path="/memories" element={<MemoryInspector />} />
            <Route path="/workflows" element={<WorkflowExecution />} />
            <Route path="/" element={<WorkflowExecution />} />
          </Routes>
        </Shell>
      </Router>
    </DataProvider>
  );
};

export default App;
