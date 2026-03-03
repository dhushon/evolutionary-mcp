import { createContext, useContext, useState, ReactNode } from 'react';

interface SelectionContextType {
  selectedId: string | null;
  setSelectedId: (id: string | null) => void;
  isCreating: boolean;
  setIsCreating: (val: boolean) => void;
}

export const SelectionContext = createContext<SelectionContextType | undefined>(undefined);

export const SelectionProvider = ({ children }: { children: ReactNode }) => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  return (
    <SelectionContext.Provider value={{ selectedId, setSelectedId, isCreating, setIsCreating }}>
      {children}
    </SelectionContext.Provider>
  );
};

export const useSelection = () => {
  const context = useContext(SelectionContext);
  if (!context) throw new Error('useSelection must be used within SelectionProvider');
  return context;
};
