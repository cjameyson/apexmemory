drop schema if exists app_code cascade;
create schema app_code;

{{ template "funcs/tg_updated_at.sql" . }}
{{ template "funcs/tg_sync_notebook_card_count.sql" . }}
