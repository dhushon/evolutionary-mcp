import React, { useContext } from 'react';
import { LogContext } from '../context/LogContext';

const LogsPanel: React.FC = () => {
    const logContext = useContext(LogContext);

    return (
        <div className="p-4 bg-gray-100 dark:bg-gray-800 shadow-md rounded-lg h-full flex flex-col">
            <h2 className="text-lg font-semibold mb-2 text-gray-900 dark:text-white">Logs Panel</h2>
            <div className="flex-grow bg-white dark:bg-gray-900 p-2 rounded-md overflow-y-auto">
                {logContext?.logs.map((log, index) => (
                    <p key={index} className="text-sm text-gray-800 dark:text-gray-300 font-mono">
                        {log}
                    </p>
                ))}
            </div>
        </div>
    );
};

export default LogsPanel;
