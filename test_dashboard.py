#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright
import sys

async def main():
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()
        
        # Capture console logs
        def console_handler(msg):
            print(f'[CONSOLE] {msg.type}: {msg.text}')
        
        page.on('console', console_handler)
        
        # Capture network requests
        def request_handler(req):
            if '/api/debug' in req.url:
                print(f'[DEBUG] {req.url}')
            if '/sse' in req.url:
                print(f'[SSE] {req.url}')
        
        page.on('request', request_handler)
        
        # Capture responses
        def response_handler(res):
            if '/sse' in res.url:
                print(f'[SSE RESPONSE] {res.status}')
        
        page.on('response', response_handler)
        
        # Navigate to dashboard
        print('Navigating to http://localhost:3001/')
        await page.goto('http://localhost:3001/', wait_until='networkidle')
        
        # Wait for Datastar to load (check for window.DatastarReady)
        print('Waiting for Datastar...')
        try:
            await page.wait_for_function('window.DatastarReady === true', timeout=10000)
            print('Datastar ready')
        except Exception as e:
            print(f'Datastar not ready: {e}')
        
        # Check if SSE connection established
        sse_connected = await page.evaluate('''() => {
            return window._sseConnection && window._sseConnection.readyState === EventSource.OPEN;
        }''')
        print(f'SSE connected: {sse_connected}')
        if sse_connected:
            print(f'SSE readyState: {await page.evaluate("window._sseConnection.readyState")}')
        
        # Wait a few seconds to see if any SSE events arrive
        print('Waiting for SSE events...')
        await asyncio.sleep(5)
        
        # Check store values
        store = await page.evaluate('''() => {
            if (window.datastar && window.datastar.store) {
                return window.datastar.store;
            }
            return null;
        }''')
        print(f'Store: {store}')
        
        # Take screenshot for debugging
        await page.screenshot(path='/tmp/dashboard.png')
        print('Screenshot saved to /tmp/dashboard.png')
        
        await browser.close()

if __name__ == '__main__':
    asyncio.run(main())