import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import DashboardLayout from './components/layout/DashboardLayout';
import WorkflowExecution from './pages/WorkflowExecution';
import GroundingManager from './pages/GroundingManager';
import { DataProvider } from './context/DataContext';

const App: React.FC = () => {
  return (
    <DataProvider>
      <Router>
        <DashboardLayout>
          <Routes>
            <Route path="/grounding" element={<GroundingManager />} />
            <Route path="/workflows" element={<WorkflowExecution />} />
            <Route path="/" element={<WorkflowExecution />} />
          </Routes>
        </DashboardLayout>
      </Router>
    </DataProvider>
  );
};

export default App;
