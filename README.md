# payments_demo
Quick golang program demonstrating simple payment transactions.

# Usage
Make sure to replace the DatabaseURL with your own firebase url.  
Also ensure that you are authenticated through default authentication with gcloud so that you can speak to firebase.

# Play
Modify the code to create cards, perform transactions, and watch your firebase database update in real time.
Firebase's built in transaction method allows for atomic writes to the database so you may spawn multiple processes to conduct thread-safe transactions.
