import pandas as pd
import mysql.connector
from mysql.connector import Error

# Step 1: Load the CSV files
def load_csv_data(complete_csv, equity_csv):
    # Load the complete.csv
    complete_df = pd.read_csv(complete_csv, compression='gzip')
    print("Complete CSV loaded:")
    print(complete_df.head())
    
    # Filter only NSE_EQ data
    complete_df = complete_df[complete_df['exchange'] == 'NSE_EQ']
    print("Filtered NSE_EQ data:")
    print(complete_df.head())

    # Extract ISIN number by splitting the instrument_key
    complete_df['ISIN'] = complete_df['instrument_key'].apply(lambda x: x.split('|')[-1])
    print("ISIN extracted from instrument_key:")
    print(complete_df[['instrument_key', 'ISIN']].head())

    # Load the equityL.csv
    equity_df = pd.read_csv(equity_csv)
    print("Equity CSV loaded:")
    print(equity_df.head())

    # Clean column names
    equity_df.columns = equity_df.columns.str.strip()
    complete_df.columns = complete_df.columns.str.strip()

    # Ensure ISIN column consistency
    equity_df['ISIN NUMBER'] = equity_df['ISIN NUMBER'].str.strip()
    complete_df['ISIN'] = complete_df['ISIN'].str.strip()

    # Merge on ISIN NUMBER and extracted ISIN
    merged_df = pd.merge(
        equity_df,
        complete_df,
        left_on='ISIN NUMBER',  # ISIN from equityL.csv
        right_on='ISIN',        # Extracted ISIN from complete.csv
        how='inner'
    )

    print("Merged DataFrame:")
    print(merged_df.head())

    # Select the required columns
    result_df = merged_df[['SYMBOL', 'instrument_key']].rename(columns={'SYMBOL': 'ticker'})
    return result_df

# Step 2: Insert data into MySQL
def insert_into_mysql(df, host, user, password, database):
    try:
        # Connect to the MySQL database
        connection = mysql.connector.connect(
            host=host,
            user=user,
            password=password,
            database=database
        )
        
        if connection.is_connected():
            print("Connected to MySQL database")

            # Prepare the SQL insert statement
            insert_query = """
                INSERT INTO TickerKeys (ticker, instrument_key)
                VALUES (%s, %s)
                ON DUPLICATE KEY UPDATE instrument_key = VALUES(instrument_key)
            """
            cursor = connection.cursor()

            # Insert each row into the database
            for _, row in df.iterrows():
                cursor.execute(insert_query, (row['ticker'], row['instrument_key']))
            
            # Commit the transaction
            connection.commit()
            print(f"Successfully inserted {len(df)} records into TickerKeys table")
    
    except Error as e:
        print(f"Error: {e}")
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()
            print("MySQL connection closed")

# Step 3: Main Execution
if __name__ == "__main__":
    COMPLETE_CSV = "./complete.csv.gz"  # Path to complete CSV (compressed)
    EQUITY_CSV = "./EQUITY_L.csv"       # Path to equityL CSV

    # MySQL database configuration
    DB_CONFIG = {
        "host": "localhost",        # Replace with your MySQL host
        "user": "root",             # Replace with your MySQL username
        "password": "root",     # Replace with your MySQL password
        "database": "stockscreener" # Replace with your database name
    }

    # Load and preprocess data
    print("Loading and preprocessing data...")
    ticker_keys_df = load_csv_data(COMPLETE_CSV, EQUITY_CSV)
    print("Data loaded successfully:")
    print(ticker_keys_df.head())

    # Insert data into MySQL
    print("Inserting data into MySQL database...")
    insert_into_mysql(ticker_keys_df, **DB_CONFIG)
