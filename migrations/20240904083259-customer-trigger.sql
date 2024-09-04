
-- +migrate Up

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION capture_customer_revision() RETURNS TRIGGER AS $$
BEGIN
    -- Insert the current row into the customer_revisions table
    INSERT INTO customer_revisions (customer_id, revision, created_at, name)
    VALUES (
        OLD.id, 
        OLD.revision, 
        OLD.created_at, 
        OLD.name
    );
    
    -- Update the revision number for the new data in the customers table
    NEW.revision := OLD.revision + 1;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER before_customer_update
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION capture_customer_revision();

-- +migrate Down

DROP TRIGGER before_customer_update ON customers;
DROP FUNCTION capture_customer_revision;