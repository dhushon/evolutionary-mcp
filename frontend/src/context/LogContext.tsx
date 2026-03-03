import { createContext, useState, ReactNode } from 'react';

interface LogContextProps {
    logs: string[];
    addLog: (message: string) => void;
}

export const LogContext = createContext<LogContextProps | undefined>(undefined);

export const LogProvider = ({ children }: { children: ReactNode }) => {
    const [logs, setLogs] = useState<string[]>([]);

    const addLog = (message: string) => {
        setLogs(prevLogs => [...prevLogs, `[${new Date().toLocaleTimeString()}] ${message}`]);
    };

    return (
        <LogContext.Provider value={{ logs, addLog }}>
            {children}
        </LogContext.Provider>
    );
};
