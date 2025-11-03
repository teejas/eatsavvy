create table if not exists public.restaurants (
    id uuid primary key default gen_random_uuid(),
    places_id varchar(255) not null,
    name text not null,
    address text,
    phone_number text not null,
    open_hours jsonb,
    nutrition_info jsonb,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
)