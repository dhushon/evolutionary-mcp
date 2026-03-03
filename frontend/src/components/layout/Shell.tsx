import React from 'react';
import Sidebar from '../Sidebar';
import Navbar from '../Navbar';
import Footer from '../Footer';
import { SelectionProvider, useSelection } from '../../context/SelectionContext';

interface ShellProps {
  children: React.ReactNode;
}

const ShellContent: React.FC<ShellProps> = ({ children }) => {
  const { selectedId, setSelectedId, setIsCreating } = useSelection();

  return (
    <div className="flex h-screen overflow-hidden bg-slate-50 dark:bg-slate-950">
      {/* Sidebar - Desktop */}
      <div className="hidden md:flex md:flex-shrink-0">
        <Sidebar 
          selectedId={selectedId} 
          onSelect={(id) => { setSelectedId(id); setIsCreating(false); }} 
          onCreate={() => { setSelectedId(null); setIsCreating(true); }} 
        />
      </div>

      <div className="flex flex-col flex-1 w-0 overflow-hidden">
        <Navbar />
        <main className="flex-1 relative overflow-y-auto focus:outline-none">
          <div className="py-6 md:py-8">
            <div className="layout-container">
              {children}
            </div>
          </div>
          <Footer />
        </main>
      </div>
    </div>
  );
};

export const Shell: React.FC<ShellProps> = ({ children }) => (
  <SelectionProvider>
    <ShellContent>{children}</ShellContent>
  </SelectionProvider>
);