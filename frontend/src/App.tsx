import React from 'react';
import AnchorManager from './components/AnchorManager';
import MemoryInspector from './components/MemoryInspector';
import AuthButton from './components/AuthButton';

const App: React.FC = () => {
  return (
    <div>
      <AuthButton />
      <AnchorManager />
      <MemoryInspector />
    </div>
  );
};

export default App;
