import { createContext, useState, ReactNode } from 'react';

interface DataContextProps {
    data: any[];
    setData: (data: any[]) => void;
}

export const DataContext = createContext<DataContextProps | undefined>(undefined);

export const DataProvider = ({ children }: { children: ReactNode }) => {
    const [data, setData] = useState<any[]>([]);

    return (
        <DataContext.Provider value={{ data, setData }}>
            {children}
        </DataContext.Provider>
    );
};
