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
            
            # Wait 10-30 seconds before next update batch
            wait_time = random.uniform(10.0, 30.0)
            await asyncio.sleep(wait_time)
            
    except KeyboardInterrupt:
        print("\nShutting down...")
    finally:
        await nc.close()

if __name__ == "__main__":
    asyncio.run(main())