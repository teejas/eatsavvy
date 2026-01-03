import { useState } from 'react';
import { Phone, MapPin, Star, Droplet, ShieldAlert, Leaf, Salad, ChevronDown, ChevronUp } from 'lucide-react';
export interface Restaurant {
  id: string;
  name: string;
  phone: string;
  address: string;
  rating: number | null;
  nutrition: {
    oil: string;
    nutFree: boolean;
    accommodations: string;
    vegetablesUsed: string[];
  };
  enrichment_status: 'pending' | 'in_progress' | 'queued' | 'completed' | 'failed';
}
interface RestaurantRowProps {
  restaurant: Restaurant;
  isSelected: boolean;
  onToggleSelect: (id: string) => void;
}
export function RestaurantRow({
  restaurant,
  isSelected,
  onToggleSelect
}: RestaurantRowProps) {
  const [expandedSection, setExpandedSection] = useState<'oil' | 'accommodations' | 'vegetables' | null>(null);

  const toggleSection = (section: 'oil' | 'accommodations' | 'vegetables') => {
    setExpandedSection(prev => prev === section ? null : section);
  };

  const statusColors: Record<Restaurant['enrichment_status'], string> = {
    completed: 'bg-green-500/10 text-green-400 border-green-500/20',
    in_progress: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20',
    queued: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
    pending: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
    failed: 'bg-red-500/10 text-red-400 border-red-500/20',
  };
  return <div className={`group relative grid grid-cols-1 md:grid-cols-12 gap-4 p-4 border-b border-sky-400/20 transition-colors duration-200 items-center ${isSelected ? 'bg-sky-500/10' : 'hover:bg-[#27272a]'}`}>
      {/* Checkbox */}
      <div className="md:col-span-1 flex items-center justify-center">
        <input type="checkbox" checked={isSelected} onChange={() => onToggleSelect(restaurant.id)} className="w-4 h-4 rounded border-zinc-600 bg-zinc-800 text-sky-500 focus:ring-sky-400 focus:ring-offset-zinc-900 cursor-pointer" aria-label={`Select ${restaurant.name}`} />
      </div>

      {/* Name Section */}
      <div className="md:col-span-3 flex flex-col justify-center">
        <h3 className="text-base font-medium text-zinc-100 group-hover:text-sky-400 transition-colors">
          {restaurant.name}
        </h3>
        <span className="text-xs text-zinc-500 font-mono mt-1 md:hidden">
          ID: {restaurant.id}
        </span>
      </div>

      {/* Contact Section */}
      <div className="md:col-span-2 flex flex-col space-y-1 text-sm text-zinc-400">
        <div className="flex items-center gap-2">
          <Phone className="w-3 h-3 text-sky-400/70" />
          <span>{restaurant.phone}</span>
        </div>
        <div className="flex items-center gap-2">
          <MapPin className="w-3 h-3 text-sky-400/70" />
          <span className="truncate" title={restaurant.address}>{restaurant.address || 'N/A'}</span>
        </div>
      </div>

      {/* Rating Section */}
      <div className="md:col-span-1 flex items-center gap-1">
        <Star className="w-4 h-4 text-amber-400" fill={restaurant.rating ? "currentColor" : "none"} />
        <span className="text-sm font-medium text-zinc-300">
          {restaurant.rating ? restaurant.rating.toFixed(1) : '—'}
        </span>
      </div>

      {/* Nutrition Section */}
      <div className="md:col-span-4 grid grid-cols-2 gap-2 text-xs">
        <div 
          className="flex flex-col bg-violet-500/5 p-2 rounded border border-violet-500/10 cursor-pointer hover:bg-violet-500/10 transition-colors"
          onClick={() => toggleSection('oil')}
        >
          <div className="flex items-center justify-between text-[10px] text-violet-300 uppercase tracking-wider mb-1">
            <div className="flex items-center gap-1">
              <Droplet className="w-3 h-3" /> Oil
            </div>
            {expandedSection === 'oil' ? (
              <ChevronUp className="w-3 h-3" />
            ) : (
              <ChevronDown className="w-3 h-3" />
            )}
          </div>
          <span className={`font-medium text-violet-400 ${expandedSection === 'oil' ? 'whitespace-pre-wrap' : 'truncate'}`}>
            {restaurant.nutrition.oil}
          </span>
        </div>
        <div className="flex flex-col bg-violet-500/5 p-2 rounded border border-violet-500/10">
          <div className="flex items-center gap-1 text-[10px] text-violet-300 uppercase tracking-wider mb-1">
            <ShieldAlert className="w-3 h-3" /> Nut Free
          </div>
          <span className="font-medium text-violet-400">
            {restaurant.nutrition.nutFree ? 'Yes' : 'No'}
          </span>
        </div>
        <div 
          className="flex flex-col bg-violet-500/5 p-2 rounded border border-violet-500/10 cursor-pointer hover:bg-violet-500/10 transition-colors"
          onClick={() => toggleSection('accommodations')}
        >
          <div className="flex items-center justify-between text-[10px] text-violet-300 uppercase tracking-wider mb-1">
            <div className="flex items-center gap-1">
              <Leaf className="w-3 h-3" /> Accommodations
            </div>
            {expandedSection === 'accommodations' ? (
              <ChevronUp className="w-3 h-3" />
            ) : (
              <ChevronDown className="w-3 h-3" />
            )}
          </div>
          <span className={`font-medium text-violet-400 ${expandedSection === 'accommodations' ? 'whitespace-pre-wrap' : 'truncate'}`}>
            {restaurant.nutrition.accommodations}
          </span>
        </div>
        <div 
          className="flex flex-col bg-violet-500/5 p-2 rounded border border-violet-500/10 cursor-pointer hover:bg-violet-500/10 transition-colors"
          onClick={() => toggleSection('vegetables')}
        >
          <div className="flex items-center justify-between text-[10px] text-violet-300 uppercase tracking-wider mb-1">
            <div className="flex items-center gap-1">
              <Salad className="w-3 h-3" /> Vegetables
            </div>
            {expandedSection === 'vegetables' ? (
              <ChevronUp className="w-3 h-3" />
            ) : (
              <ChevronDown className="w-3 h-3" />
            )}
          </div>
          <span className={`font-medium text-violet-400 ${expandedSection === 'vegetables' ? 'whitespace-pre-wrap' : 'truncate'}`}>
            {restaurant.nutrition.vegetablesUsed.join(', ')}
          </span>
        </div>
      </div>

      {/* Status Section */}
      <div className="md:col-span-1 flex justify-start md:justify-end">
        <span className={`px-2 py-1 rounded-full text-[10px] uppercase font-bold tracking-wider border ${statusColors[restaurant.enrichment_status]}`}>
          {restaurant.enrichment_status}
        </span>
      </div>

      {/* Mobile-only divider */}
      <div className="absolute bottom-0 left-4 right-4 h-px bg-sky-400/10 md:hidden" />
    </div>;
}