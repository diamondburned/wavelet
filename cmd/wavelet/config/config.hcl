wallet = "config/wallet.txt"
port = 3000
host = "127.0.0.1"

api {
  port = 0
}

system {
  // Timeout for querying a transaction to K peers.
  // In seconds.
  query_timeout = 10

  // Max graph depth difference to search for eligible transaction
  // parents from for our node.
  max_eligible_parents_depth_diff = 5

  // Number of ancestors to derive a median timestamp from.
  median_timestamp_num_ancestors = 5

  validator_reward_amount = 2

  // In milliseconds.
  expected_consensus_time = 1000

  critical_timestamp_average_window_size = 3

  min_stake = 100

  // Snowball consensus protocol parameters.
  snowball {
    k = 1
    alpha = 0.8
    beta = 10
  }

  // Difficulty to define a critical transaction.
  difficulty {
    min = 5
    max = 16
  }
}
