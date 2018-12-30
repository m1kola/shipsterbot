BEGIN;

-- `Add` corresponds to the commandAdd constant from bot/telegram/commands.go
UPDATE unfinished_commands SET command='ADD_SHOPPING_ITEM' WHERE command='add';

COMMIT;
