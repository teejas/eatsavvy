import { useState, useMemo } from 'react';
import { Phone, MapPin, Star, Droplet, ShieldAlert, Leaf, Salad, ChevronDown, ChevronUp, Clock } from 'lucide-react';

// TimePoint represents a specific time on a weekday (from API, in UTC)
export interface TimePoint {
  weekday: number; // 0 = Sunday, 1 = Monday, etc.
  hour: number;
  minute: number;
}

// TimeRange represents an open/close time range (from API, in UTC)
export interface TimeRange {
  open: TimePoint;
  close: TimePoint;
}

export interface Restaurant {
  id: string;
  name: string;
  phone: string;
  address: string;
  rating: number | null;
  openHours: TimeRange[] | null;
  nutrition: {
    oil: string;
    nutFree: boolean;
    accommodations: string;
    vegetablesUsed: string;
  };
  enrichment_status: 'pending' | 'in_progress' | 'queued' | 'completed' | 'failed';
}

const WEEKDAY_NAMES = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

// Convert UTC time to local time, returning the adjusted weekday, hour, and minute
function utcToLocal(weekday: number, hour: number, minute: number): { weekday: number; hour: number; minute: number } {
  // Get the timezone offset in minutes (negative for west of UTC, e.g., PST is -480)
  const offsetMinutes = new Date().getTimezoneOffset();
  
  // Create a date object for the given UTC time (using a reference week)
  // We use Jan 5, 2025 as reference (it's a Sunday)
  const referenceDate = new Date(Date.UTC(2025, 0, 5 + weekday, hour, minute));
  
  // Apply the timezone offset
  const localDate = new Date(referenceDate.getTime() - offsetMinutes * 60 * 1000);
  
  return {
    weekday: localDate.getUTCDay(),
    hour: localDate.getUTCHours(),
    minute: localDate.getUTCMinutes(),
  };
}

// Format time as 12-hour with AM/PM
function formatTime(hour: number, minute: number): string {
  const period = hour >= 12 ? 'PM' : 'AM';
  const displayHour = hour % 12 || 12;
  const displayMinute = minute.toString().padStart(2, '0');
  return `${displayHour}:${displayMinute} ${period}`;
}

// Group time ranges by weekday (after converting to local time)
function groupHoursByWeekday(openHours: TimeRange[]): Map<number, { open: string; close: string; openMinutes: number; closeMinutes: number }[]> {
  const grouped = new Map<number, { open: string; close: string; openMinutes: number; closeMinutes: number }[]>();
  
  for (const range of openHours) {
    const localOpen = utcToLocal(range.open.weekday, range.open.hour, range.open.minute);
    const localClose = utcToLocal(range.close.weekday, range.close.hour, range.close.minute);
    
    const weekday = localOpen.weekday;
    if (!grouped.has(weekday)) {
      grouped.set(weekday, []);
    }
    grouped.get(weekday)!.push({
      open: formatTime(localOpen.hour, localOpen.minute),
      close: formatTime(localClose.hour, localClose.minute),
      openMinutes: localOpen.hour * 60 + localOpen.minute,
      closeMinutes: localClose.hour * 60 + localClose.minute,
    });
  }
  
  return grouped;
}

// Check if the restaurant is currently open based on grouped hours
function isCurrentlyOpen(groupedHours: Map<number, { open: string; close: string; openMinutes: number; closeMinutes: number }[]> | null): boolean {
  if (!groupedHours) return false;
  
  const now = new Date();
  const currentWeekday = now.getDay();
  const currentMinutes = now.getHours() * 60 + now.getMinutes();
  
  const todayHours = groupedHours.get(currentWeekday);
  if (!todayHours) return false;
  
  return todayHours.some(({ openMinutes, closeMinutes }) => {
    // Handle overnight hours (e.g., 10 PM - 2 AM)
    if (closeMinutes < openMinutes) {
      return currentMinutes >= openMinutes || currentMinutes < closeMinutes;
    }
    return currentMinutes >= openMinutes && currentMinutes < closeMinutes;
  });
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
  const [expandedSection, setExpandedSection] = useState<'oil' | 'accommodations' | 'vegetables' | 'hours' | null>(null);

  const toggleSection = (section: 'oil' | 'accommodations' | 'vegetables' | 'hours') => {
    setExpandedSection(prev => prev === section ? null : section);
  };

  // Memoize the grouped hours to avoid recalculating on every render
  const groupedHours = useMemo(() => {
    if (!restaurant.openHours || restaurant.openHours.length === 0) return null;
    return groupHoursByWeekday(restaurant.openHours);
  }, [restaurant.openHours]);

  // Check if currently open
  const isOpen = useMemo(() => isCurrentlyOpen(groupedHours), [groupedHours]);

  const statusColors: Record<Restaurant['enrichment_status'], string> = {
    completed: 'bg-green-500/10 text-green-400 border-green-500/20',
    in_progress: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20',
    queued: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
    pending: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
    failed: 'bg-red-500/10 text-red-400 border-red-500/20',
  };

  const displayStatus = restaurant.enrichment_status || 'pending';
  return <div className={`group relative grid grid-cols-1 md:grid-cols-12 gap-4 p-4 border-b border-sky-400/20 transition-colors duration-200 items-start ${isSelected ? 'bg-sky-500/10' : 'hover:bg-[#27272a]'}`}>
      {/* Checkbox */}
      <div className="md:col-span-1 flex items-center justify-center pt-1">
        <input type="checkbox" checked={isSelected} onChange={() => onToggleSelect(restaurant.id)} className="w-4 h-4 rounded border-zinc-600 bg-zinc-800 text-sky-500 focus:ring-sky-400 focus:ring-offset-zinc-900 cursor-pointer" aria-label={`Select ${restaurant.name}`} />
      </div>

      {/* Name Section */}
      <div className="md:col-span-2 flex flex-col justify-center">
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
        <div className="flex items-start gap-2">
          <MapPin className="w-3 h-3 text-sky-400/70 mt-0.5 flex-shrink-0" />
          <span className="whitespace-pre-line">{restaurant.address ? restaurant.address.split(', ').join(',\n') : 'N/A'}</span>
        </div>
      </div>

      {/* Rating Section */}
      <div className="md:col-span-1 flex items-center gap-1">
        <Star className="w-4 h-4 text-amber-400" fill={restaurant.rating ? "currentColor" : "none"} />
        <span className="text-sm font-medium text-zinc-300">
          {restaurant.rating ? restaurant.rating.toFixed(1) : '—'}
        </span>
      </div>

      {/* Open Hours Section */}
      <div className="md:col-span-2">
        <div 
          className={`flex flex-col p-2 rounded border cursor-pointer transition-colors ${
            isOpen 
              ? 'bg-emerald-500/5 border-emerald-500/10 hover:bg-emerald-500/10' 
              : 'bg-zinc-500/5 border-zinc-500/10 hover:bg-zinc-500/10'
          }`}
          onClick={() => toggleSection('hours')}
        >
          <div className={`flex items-center justify-between text-[10px] uppercase tracking-wider mb-1 ${
            isOpen ? 'text-emerald-300' : 'text-zinc-400'
          }`}>
            <div className="flex items-center gap-1">
              <Clock className="w-3 h-3" /> {isOpen ? 'Open Now' : 'Closed'}
            </div>
            {expandedSection === 'hours' ? (
              <ChevronUp className="w-3 h-3" />
            ) : (
              <ChevronDown className="w-3 h-3" />
            )}
          </div>
          {groupedHours ? (
            expandedSection === 'hours' ? (
              <div className="space-y-0.5">
                {WEEKDAY_NAMES.map((name, idx) => {
                  const hours = groupedHours.get(idx);
                  const isToday = idx === new Date().getDay();
                  return (
                    <div key={idx} className={`flex justify-between text-xs ${isToday ? 'font-semibold' : ''}`}>
                      <span className={`w-8 ${isOpen && isToday ? 'text-emerald-400' : isToday ? 'text-zinc-300' : 'text-zinc-500'}`}>
                        {name}
                      </span>
                      <span className={isOpen && isToday ? 'text-emerald-300/80' : 'text-zinc-400'}>
                        {hours ? hours.map(h => `${h.open}–${h.close}`).join(', ') : 'Closed'}
                      </span>
                    </div>
                  );
                })}
              </div>
            ) : (
              <span className={`font-medium text-xs ${isOpen ? 'text-emerald-400' : 'text-zinc-500'}`}>
                {groupedHours.size} days · Click to expand
              </span>
            )
          ) : (
            <span className="font-medium text-zinc-500 text-xs">Not available</span>
          )}
        </div>
      </div>

      {/* Nutrition Section */}
      <div className="md:col-span-3 grid grid-cols-2 gap-2 text-xs">
        <div 
          className="flex flex-col bg-violet-500/5 p-2 rounded border border-violet-500/10 cursor-pointer hover:bg-violet-500/10 transition-colors"
          onClick={() => toggleSection('oil')}
        >
          <div className="flex items-center justify-between text-[10px] text-violet-300 uppercase tracking-wider mb-1">
            <div className="flex items-center gap-1">
              <Droplet className="w-3 h-3" /> Cooking Oil
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
              <Salad className="w-3 h-3" /> Vegetables Used
            </div>
            {expandedSection === 'vegetables' ? (
              <ChevronUp className="w-3 h-3" />
            ) : (
              <ChevronDown className="w-3 h-3" />
            )}
          </div>
          <span className={`font-medium text-violet-400 ${expandedSection === 'vegetables' ? 'whitespace-pre-wrap' : 'truncate'}`}>
            {restaurant.nutrition.vegetablesUsed}
          </span>
        </div>
      </div>

      {/* Status Section */}
      <div className="md:col-span-1 flex justify-start md:justify-end">
        <span className={`px-2 py-1 rounded-full text-[10px] uppercase font-bold tracking-wider border ${statusColors[displayStatus]}`}>
          {displayStatus}
        </span>
      </div>

      {/* Mobile-only divider */}
      <div className="absolute bottom-0 left-4 right-4 h-px bg-sky-400/10 md:hidden" />
    </div>;
}