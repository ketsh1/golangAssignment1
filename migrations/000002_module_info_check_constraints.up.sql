ALTER TABLE module_info ADD CONSTRAINT modules_duration_check CHECK (module_duration > 5 AND module_duration <=15);
ALTER TABLE module_info ADD CONSTRAINT check_dates CHECK (updated_at >= created_at);
