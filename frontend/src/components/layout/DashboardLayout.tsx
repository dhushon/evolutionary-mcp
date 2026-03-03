import React, { useState, createContext, useContext } from 'react';
import Navbar from '../Navbar';
import Sidebar from '../Sidebar';
import Footer from '../Footer';

interface SelectionContextType {
  selectedId: string | null;
  setSelectedId: (id: string | null) => void;
  isCreating: boolean;
  setIsCreating: (val: boolean) => void;
}

export const SelectionContext = createContext<SelectionContextType | undefined>(undefined);

export const useSelection = () => {
  const context = useContext(SelectionContext);
  if (!context) throw new Error('useSelection must be used within SelectionProvider');
  return context;
};

interface DashboardLayoutProps {
  children: React.ReactNode;
}

const DashboardLayout: React.FC<DashboardLayoutProps> = ({ children }) => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  const handleSelect = (id: string) => {
    setSelectedId(id);
    setIsCreating(false);
  };

  const handleCreate = () => {
    setSelectedId(null);
    setIsCreating(true);
  };

  return (
    <SelectionContext.Provider value={{ selectedId, setSelectedId, isCreating, setIsCreating }}>
      <div className="flex flex-col h-screen bg-bg-base text-text-base font-sans transition-colors duration-300">
        <Navbar />
        <div className="flex flex-1 overflow-hidden">
          <Sidebar
            selectedId={selectedId}
            onSelect={handleSelect}
            onCreate={handleCreate}
          />
          <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
            <main className="flex-1 overflow-y-auto">
              {children}
            </main>
            <Footer />
          </div>
        </div>
      </div>
    </SelectionContext.Provider>
  );
};

export default DashboardLayout;
