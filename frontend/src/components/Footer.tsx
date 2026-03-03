export default function Footer() {
    return (
        <footer className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 py-3 px-6">
            <div className="flex justify-between items-center text-xs text-gray-500 dark:text-gray-400">
                <p>&copy; {new Date().getFullYear()} Evolutionary Memory MCP. All rights reserved.</p>
                <div className="flex space-x-4">
                    <span>v0.1.0-alpha</span>
                </div>
            </div>
        </footer>
    );
}