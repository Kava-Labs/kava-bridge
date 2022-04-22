# Gas Estimation and Fees

Every peer in a signing party has to estimate the gas with some margin, in
addition to the fee to deduct from the transferred amount.

As Ethereum gas fees are high and Kava fees are low, fees for transactions going
into Kava are paid by the relayer and fees for outgoing transactions to Ethereum
are paid by deducting the corresponding gas fee amount from the amount being
transferred.

After each peer estimates gas and determines the amount to deduct from the
transferred amount, all peers will need to agree on the same value to sign. This
is done by broadcasting peer gas and fee values and each peer picking the
median.
