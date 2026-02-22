#!/usr/bin/env python3
import asyncio
import json
import random
import time
import sys
import nats

async def main():
    # Connect to NATS
    try:
        print("Attempting to connect to NATS...")
        nc = await nats.connect("nats://nats:4222")
        print("Connected to NATS")
    except Exception as e:
        print(f"Failed to connect to NATS: {e}")
        sys.exit(1)
    
    # Publish random numbers to Quarter A topics
    try:
        while True:
            # Publish to Quarter A topics (24 topics)
            for col in range(1, 4):  # Columns 1-3
                for num in range(1, 9):  # Numbers 1-8
                    topic = f"quarterA.col{col}.num{num}"
                    value = random.randint(1, 99)
                    
                    # Publish as JSON number
                    await nc.publish(topic, json.dumps(value).encode())
                    print(f"Published to {topic}: {value}")
            
            # Publish to Quarter D topics (6 topics)
            for col in range(1, 7):  # Columns 1-6
                topic = f"quarterD.col{col}"
                value = random.randint(1, 99)
                
                # Publish as JSON number
                await nc.publish(topic, json.dumps(value).encode())
                print(f"Published to {topic}: {value}")
            print('DEBUG: After Quarter D')
            # Publish to Quarter C topics (wait times)
            print("DEBUG: Publishing Quarter C wait times")
            for cbw in range(1, 4):
                print(f"DEBUG: Iteration cbw={cbw}")
                topic = f"quarterC.cbw{cbw}.waitTime"
                value = random.randint(1, 30)  # wait time minutes
                print(f"DEBUG: Publishing {topic}={value}")
                await nc.publish(topic, json.dumps(value).encode())
                print(f"Published to {topic}: {value}")
            
            # Publish to Quarter C wasted minutes topics
            print("DEBUG: Publishing Quarter C wasted minutes")
            wasted_keys = ["hour4", "hour3", "hour2", "hour1", "current"]
            for key in wasted_keys:
                print(f"DEBUG: Wasted key={key}")
                topic = f"quarterC.wastedMinutes.{key}"
                value = random.randint(50, 200)  # wasted minutes
                print(f"DEBUG: Publishing {topic}={value}")
                await nc.publish(topic, json.dumps(value).encode())
                print(f"Published to {topic}: {value}")
            
            # Wait 10-30 seconds before next update batch
            print("DEBUG: Sleeping before next batch")
            wait_time = random.uniform(10.0, 30.0)
            await asyncio.sleep(wait_time)
            
    except KeyboardInterrupt:
        print("\nShutting down...")
    finally:
        await nc.close()

if __name__ == "__main__":
    asyncio.run(main())