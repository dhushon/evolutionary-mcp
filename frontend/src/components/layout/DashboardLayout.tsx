import React from 'react';
import { SelectionProvider, useSelection } from '../../context/SelectionContext';
import Footer from '../Footer';
import Navbar from '../Navbar';
import Sidebar from '../Sidebar';

interface DashboardLayoutProps {
  children: React.ReactNode;
}

const LayoutContent: React.FC<DashboardLayoutProps> = ({ children }) => {
  const { selectedId, setSelectedId, setIsCreating } = useSelection();

  const handleSelect = (id: string) => {
    setSelectedId(id);
    setIsCreating(false);
  };

  const handleCreate = () => {
    setSelectedId(null);
    setIsCreating(true);
  };

  return (
    <div className="grid grid-rows-[auto_1fr] h-screen w-full bg-bg-base text-text-base">
      {/* Row 1: Navbar */}
      <div className="border-b border-border-base">
        <Navbar />
      </div>

      {/* Row 2: Content */}
      <div className="grid grid-cols-[auto_1fr] overflow-hidden">
        {/* Sidebar */}
        <div className="border-r border-border-base overflow-y-auto">
          <Sidebar
            selectedId={selectedId}
            onSelect={handleSelect}
            onCreate={handleCreate}
          />
        </div>

        {/* Main Area */}
        <div className="flex flex-col overflow-hidden">
          <main className="flex-1 overflow-y-auto p-4">
            {children}
          </main>
          <div className="border-t border-border-base">
            <Footer />
          </div>
        </div>
      </div>
    </div>
  );
};

const DashboardLayout: React.FC<DashboardLayoutProps> = ({ children }) => {
  return (
    <SelectionProvider>
      <LayoutContent>{children}</LayoutContent>
    </SelectionProvider>
  );
};

export default DashboardLayout;
