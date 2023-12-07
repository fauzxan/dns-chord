'''
Code to analyze the log outputs.
This code will help analyze two things:
1. Hop Count - indicating the number nodes the message hopped over for "n" queries
2. Time passage - amount of time it takes for the system to return a successfull result for "n" queries
'''

import os
import json
import time

def process_logs():
    hop_count = {} 
    time_list = {}
    for file in os.listdir("./"):
        if ".log" not in file:
            continue
        print("Reading", file)
        file_path = os.path.join("./", file)
        with open(file_path) as f:
            try:
                for line in f.readlines():
                    json_content = json.loads(line)
                    if json_content:
                        if "HOP COUNT" in json_content["message"]: 
                            # Code to dump hop count into hop count dictionary/list
                            print(json_content["message"])
                        elif "TIME" in json_content["message"]:
                            # code to dump time messages in time dictionary/list
                            print(json_content["message"])

            except json.JSONDecodeError as e:
                print("ERROR", e)
    return hop_count, time_list

def main():
    print("Processing logs now...")
    start = time.time()
    process_logs()
    print("Done in", time.time() - start)



if __name__ == "__main__":
    main()