import { useMemo, useState, useEffect } from 'react';
import { SearchBar } from './SearchBar';
import { RestaurantRow, Restaurant, TimeRange } from './RestaurantRow';
import { UtensilsCrossed, Sparkles, Loader2, AlertCircle } from 'lucide-react';

const API_BASE_URL = import.meta.env.VITE_EATSAVVY_API_URL || 'https://api.eatsavvy.org';
const API_KEY = import.meta.env.VITE_EATSAVVY_API_KEY;

// Helper to create authenticated fetch requests
function authFetch(url: string, options: RequestInit = {}): Promise<Response> {
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      ...(API_KEY ? { 'Authorization': `Bearer ${API_KEY}` } : {}),
    },
  });
}

// API response types
interface ApiNutritionInfo {
  oil: string;
  nutFree: boolean;
  accommodations: string;
  vegetables: string;
}

interface ApiRestaurant {
  id: string;
  name: string;
  address: string;
  phoneNumber: string;
  openHours: TimeRange[] | null;
  nutritionInfo: ApiNutritionInfo | null;
  rating: number | null;
  enrichmentStatus: Restaurant['enrichment_status'];
}

// Transform API response to frontend Restaurant type
function transformRestaurant(api: ApiRestaurant): Restaurant {
  return {
    id: api.id,
    name: api.name,
    phone: api.phoneNumber || 'N/A',
    address: api.address || '',
    rating: api.rating,
    openHours: api.openHours,
    nutrition: {
      oil: api.nutritionInfo?.oil || 'Unknown',
      nutFree: api.nutritionInfo?.nutFree || false,
      accommodations: api.nutritionInfo?.accommodations || 'None',
      vegetablesUsed: api.nutritionInfo?.vegetables || 'Unknown',
    },
    enrichment_status: api.enrichmentStatus,
  };
}
export function RestaurantListView() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [restaurants, setRestaurants] = useState<Restaurant[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchRestaurants() {
      try {
        setLoading(true);
        setError(null);
        
        const response = await authFetch(`${API_BASE_URL}/restaurant`);
        if (!response.ok) {
          throw new Error(`Failed to fetch restaurants: ${response.statusText}`);
        }
        const apiRestaurants: ApiRestaurant[] = await response.json();
        setRestaurants(apiRestaurants.map(transformRestaurant));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unexpected error occurred');
      } finally {
        setLoading(false);
      }
    }

    fetchRestaurants();
  }, []);

  // Local filtering for real-time search as you type
  const filteredData = useMemo(() => {
    const query = searchQuery.toLowerCase();
    return restaurants.filter(item => 
      item.name.toLowerCase().includes(query) || 
      item.phone.includes(query) || 
      item.address.toLowerCase().includes(query)
    );
  }, [searchQuery, restaurants]);

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      // If search is empty, reload all restaurants
      try {
        setLoading(true);
        setError(null);
        const response = await authFetch(`${API_BASE_URL}/restaurant`);
        if (!response.ok) {
          throw new Error(`Failed to fetch restaurants: ${response.statusText}`);
        }
        const apiRestaurants: ApiRestaurant[] = await response.json();
        setRestaurants(apiRestaurants.map(transformRestaurant));
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unexpected error occurred');
      } finally {
        setLoading(false);
      }
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await authFetch(`${API_BASE_URL}/search`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ query: searchQuery }),
      });
      if (!response.ok) {
        throw new Error(`Search failed: ${response.statusText}`);
      }
      const apiRestaurants: ApiRestaurant[] = await response.json();
      setRestaurants(apiRestaurants.map(transformRestaurant));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };
  const handleToggleSelect = (id: string) => {
    setSelectedIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };
  const handleSelectAll = () => {
    if (selectedIds.size === filteredData.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(filteredData.map(r => r.id)));
    }
  };
  const [enriching, setEnriching] = useState(false);

  const handleEnrich = async () => {
    const ids = Array.from(selectedIds);
    if (ids.length === 0) return;

    try {
      setEnriching(true);
      setError(null);
      const response = await authFetch(`${API_BASE_URL}/enrich`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ ids }),
      });
      if (!response.ok) {
        throw new Error(`Enrich failed: ${response.statusText}`);
      }
      const enrichedRestaurants: ApiRestaurant[] = await response.json();
      
      // Update the restaurants list with enriched data
      setRestaurants(prev => {
        const enrichedMap = new Map(enrichedRestaurants.map(r => [r.id, transformRestaurant(r)]));
        return prev.map(r => enrichedMap.get(r.id) || r);
      });
      
      // Clear selection after successful enrichment
      setSelectedIds(new Set());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to enrich restaurants');
    } finally {
      setEnriching(false);
    }
  };
  const allSelected = filteredData.length > 0 && selectedIds.size === filteredData.length;
  return <div className="min-h-screen bg-zinc-900 text-zinc-100 font-sans selection:bg-sky-400/30 selection:text-sky-200">
      {/* Sticky Header */}
      <header className="sticky top-0 z-10 bg-zinc-900/95 backdrop-blur-sm border-b border-sky-400/30 shadow-lg shadow-black/20">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-sky-400/10 rounded-lg border border-sky-400/20">
                <UtensilsCrossed className="w-6 h-6 text-sky-400" />
              </div>
              <div>
                <h1 className="text-xl font-bold text-zinc-100 tracking-tight">
                  Restaurant Directory
                </h1>
                <p className="text-xs text-zinc-400">
                  {filteredData.length}{' '}
                  {filteredData.length === 1 ? 'result' : 'results'} found
                  {selectedIds.size > 0 && ` · ${selectedIds.size} selected`}
                </p>
              </div>
            </div>
            <div className="w-full md:w-auto">
              <SearchBar value={searchQuery} onChange={setSearchQuery} onSearch={handleSearch} />
            </div>
          </div>

          {/* Column Headers (Desktop) */}
          <div className="hidden md:grid grid-cols-12 gap-4 mt-6 px-4 pb-2 text-xs font-semibold text-sky-400 uppercase tracking-wider opacity-80">
            <div className="col-span-1 flex items-center justify-center">
              <input type="checkbox" checked={allSelected} onChange={handleSelectAll} className="w-4 h-4 rounded border-zinc-600 bg-zinc-800 text-sky-500 focus:ring-sky-400 focus:ring-offset-zinc-900 cursor-pointer" aria-label="Select all" />
            </div>
            <div className="col-span-2">Restaurant Details</div>
            <div className="col-span-2">Contact Info</div>
            <div className="col-span-1">Rating</div>
            <div className="col-span-2">Open Hours</div>
            <div className="col-span-3">Nutrition Info</div>
            <div className="col-span-1 text-right">Status</div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 pb-20">
        <div className="bg-zinc-900 border-x border-sky-400/10 min-h-[500px]">
          {loading ? (
            <div className="flex flex-col items-center justify-center py-20 text-zinc-400">
              <Loader2 className="w-12 h-12 mb-4 animate-spin text-sky-400" />
              <p className="text-lg font-medium">Loading restaurants...</p>
            </div>
          ) : error ? (
            <div className="flex flex-col items-center justify-center py-20 text-red-400">
              <AlertCircle className="w-12 h-12 mb-4" />
              <p className="text-lg font-medium">Error loading restaurants</p>
              <p className="text-sm text-zinc-500 mt-1">{error}</p>
              <button 
                onClick={() => window.location.reload()} 
                className="mt-4 px-4 py-2 bg-sky-500/20 text-sky-400 rounded-lg hover:bg-sky-500/30 transition-colors"
              >
                Try Again
              </button>
            </div>
          ) : filteredData.length > 0 ? (
            <div className="divide-y divide-sky-400/10">
              {filteredData.map(restaurant => (
                <RestaurantRow 
                  key={restaurant.id} 
                  restaurant={restaurant} 
                  isSelected={selectedIds.has(restaurant.id)} 
                  onToggleSelect={handleToggleSelect} 
                />
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-20 text-zinc-500">
              <SearchIcon className="w-12 h-12 mb-4 opacity-20" />
              <p className="text-lg font-medium">No restaurants found</p>
              <p className="text-sm">Try adjusting your search terms</p>
            </div>
          )}
        </div>
      </main>

      {/* Floating Enrich Button */}
      {selectedIds.size > 0 && <div className="fixed bottom-8 left-1/2 transform -translate-x-1/2 z-20 animate-in slide-in-from-bottom-4 duration-300">
          <button 
            onClick={handleEnrich} 
            disabled={enriching}
            className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-violet-600 hover:from-indigo-500 hover:to-violet-500 disabled:from-indigo-800 disabled:to-violet-800 disabled:cursor-not-allowed text-white rounded-full font-semibold shadow-lg shadow-indigo-500/50 transition-all duration-200 hover:shadow-xl hover:shadow-indigo-500/60 hover:scale-105 disabled:hover:scale-100"
          >
            {enriching ? (
              <Loader2 className="w-5 h-5 animate-spin" />
            ) : (
              <Sparkles className="w-5 h-5" />
            )}
            <span>
              {enriching ? 'Enriching...' : `Enrich ${selectedIds.size} ${selectedIds.size === 1 ? 'Restaurant' : 'Restaurants'}`}
            </span>
          </button>
        </div>}
    </div>;
}
function SearchIcon({
  className
}: {
  className?: string;
}) {
  return <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
      <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
    </svg>;
}