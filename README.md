# health-authority
The user's health data contains sensitive information that needs to be kept private. We need to utilize zero-knowledge proofs to ensure that sensitive data remains encrypted and under the control of the data owner at all times. Aleo ZKML addresses this issue by integrating privacy into the smart contract layer. For third parties that require verification of user data, ZKP (Zero-Knowledge Proof) smart contracts can be written to prove the correctness of computations without revealing the underlying data.
# Program Description
## Program URLs:
	* 1. https://github.com/scaleomngt/health-authority/tree/main/front/health_data --Frontend code
	* 2. https://github.com/scaleomngt/health-authority --Server code
	* 3. https://github.com/scaleomngt/health-authority/tree/main/aleo/health --Leo code
# Installation Instructions
## Frontend Deployment
(1) Step 1: Download the dependencies required for the program
```
npm install
```
(2) Step 2: Verify if all dependencies are successfully downloaded by running the program locally
```
npm run serve  
```
(3) Step 3: Access the configured IP address and port in a browser to use the application

## Backend Deployment
```
go build
```
## leoDeployment
```
- `cd aleo/health && leo build`
- `export PRIVATE_KEY=<PRIVATE_KEY>`
- Publish£º
snarkos developer deploy health.aleo --private-key $PRIVATE_KEY --query "https://vm.aleo.org/api" --path "build/" \
--broadcast "https://vm.aleo.org/api/testnet3/transaction/broadcast" --fee 600000 \
--record <Record_that_you_just_transferred_credits_to>
```
# Data flow diagram
<img src="https://github.com/scaleomngt/health-authority/blob/main/t5.png" alt="drawing" width="800"/>

# Program Execution
* 1. To register an Aleo account for users on a website page.
* 2. Enter user health data, such as blood pressure, heart rate, blood sugar, and other information, then click submit.
* 3. Invoke the Leo verification program and receive the transaction hash.
* 4. Use the transaction hash to provide third parties with confirmation of the user's health status without exposing the user's private data.
<img src="https://github.com/scaleomngt/health-authority/blob/main/t1.png" alt="drawing" width="800"/>
<img src="https://github.com/scaleomngt/health-authority/blob/main/t2.png" alt="drawing" width="800"/>
<img src="https://github.com/scaleomngt/health-authority/blob/main/t3.png" alt="drawing" width="800"/>

