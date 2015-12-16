## Heavy Duty
1. Build Maps - Year, Make, Model, Style
2. Build "New" Temp Tables
3. Read CSV 
4. For each line:
  1. User Prompt - Ask user for make choice (e.g. Chevy/GMC)
  2. Parse into multiple VehicleParts

5. Enter VehicleParts into "New" Database Tables
  1. Missing Year, Make, Model, Style?
    1. User Prompt - Ask user if they would like to enter it 
    	1. Enter in DB (temp tables) & map -or-
    	2. Enter VehiclePart into Error Output file
  2. Enter Vehicle into Vehicle Temp table if it doesn't exist already
  3. Enter VehiclePart into VehiclePart Temp table if it doesn't exist already

6. User Prompt - Merge "New" DB tables (Vehicle, VehiclePart) into existing tables (after user inspection)
7. Index all these parts (mongo/redis/iapi/goqueue)