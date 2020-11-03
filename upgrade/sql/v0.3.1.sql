ALTER TABLE public.express_provider
    ALTER COLUMN enabled TYPE int2 USING enabled::int;