CREATE OR REPLACE FUNCTION app_code.tg_sync_notebook_card_count()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE app.notebooks
        SET total_cards = total_cards + 1, updated_at = now()
        WHERE user_id = NEW.user_id AND id = NEW.notebook_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE app.notebooks
        SET total_cards = total_cards - 1, updated_at = now()
        WHERE user_id = OLD.user_id AND id = OLD.notebook_id;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' AND OLD.notebook_id IS DISTINCT FROM NEW.notebook_id THEN
        UPDATE app.notebooks
        SET total_cards = total_cards - 1, updated_at = now()
        WHERE user_id = OLD.user_id AND id = OLD.notebook_id;
        UPDATE app.notebooks
        SET total_cards = total_cards + 1, updated_at = now()
        WHERE user_id = NEW.user_id AND id = NEW.notebook_id;
        RETURN NEW;
    END IF;
    RETURN NULL;
END$$;
