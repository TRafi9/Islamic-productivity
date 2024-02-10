-- USER ENTRIES FOR SUBMITTING YES/NO ANSWER :

	-- // insertSQL := `
	-- // CREATE TABLE IF NOT EXISTS user_submissions (
	-- // 	random_primary_key UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	-- // 	user_id VARCHAR(255),
	-- // 	productive_val BOOLEAN,
	-- // 	first_prayer_name VARCHAR(255),
	-- // 	second_prayer_name VARCHAR(255),
	-- // 	first_prayer_time TIMESTAMP,
	-- // 	second_prayer_time TIMESTAMP,
	-- // 	ingestion_timestamp TIMESTAMP
	-- // );
	-- // `





-- CREATING USER STATEMENTS:
	-- // primary key is users email here to stop multiple registrations of same user email going through
	-- // createTableSQL := `
	-- 	// CREATE TABLE IF NOT EXISTS users (
	-- 	// 	user_id VARCHAR(255) PRIMARY KEY,
	-- 	// 	password_hash VARCHAR(255),
	-- 	// 	creation_timestamp TIMESTAMP
	-- 	//  verified_email BOOLEAN
	-- 	// );
	-- 	// `





-- VERIFY USER REGISTRATION STATEMENTS:
	-- // // setup verification table in db to hold verification codes for users to register
	-- // createTableSQL := `
	-- // CREATE TABLE IF NOT EXISTS email_verification_check (
	-- // 	user_id VARCHAR(255) PRIMARY KEY,
	-- // 	email_verification_code INT,
	-- // 	expiry_time TIMESTAMP
	-- // );


