import React from 'react';
import { Search } from 'lucide-react';
interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  onSearch: () => void;
}
export function SearchBar({
  value,
  onChange,
  onSearch
}: SearchBarProps) {
  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      onSearch();
    }
  };
  return <div className="relative w-full flex gap-2">
      <div className="relative flex-1">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-zinc-400" />
        </div>
        <input type="text" value={value} onChange={e => onChange(e.target.value)} onKeyPress={handleKeyPress} className="block w-full pl-10 pr-3 py-2 border border-zinc-700 rounded-md leading-5 bg-zinc-800 text-zinc-100 placeholder-zinc-500 focus:outline-none focus:ring-1 focus:ring-sky-400 focus:border-sky-400 sm:text-sm transition-colors duration-200" placeholder="Search restaurants..." aria-label="Search restaurants" />
      </div>
    </div>;
}