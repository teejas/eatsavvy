create table if not exists public.calls (
    id uuid primary key default gen_random_uuid(),
    places_id varchar(255) not null unique,
    vapi_call_id text,
    call_status text,
    transcript text,
    structured_outputs jsonb,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
)