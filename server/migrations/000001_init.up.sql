CREATE TABLE IF NOT EXISTS public.order_items
(
    id serial NOT NULL,
    order_id integer NOT NULL,
    product_id integer NOT NULL,
    CONSTRAINT order_items_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.orders
(
    id serial NOT NULL,
    cart_id integer NOT NULL,
    create_at timestamp without time zone NOT NULL DEFAULT now(),
    delivered_at timestamp without time zone NOT NULL,
    transport_type text COLLATE pg_catalog."default" NOT NULL,
    user_id integer NOT NULL,
    address text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.order_items
    ADD CONSTRAINT fk_order_items FOREIGN KEY (order_id)
        REFERENCES public.orders (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE;