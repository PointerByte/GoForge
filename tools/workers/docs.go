// Package workers provides a simple in-process task dispatcher backed by a
// buffered function queue.
//
// The package exposes a process-wide worker loop. Tasks are submitted with
// AddTask and executed asynchronously after RunWorkers starts the dispatcher.
// SetWorkersLimit replaces the queue with a new buffered channel whose capacity
// controls how many tasks can be queued before AddTask blocks.
//
// Main entry points:
//   - SetWorkersLimit configures the task queue capacity.
//   - AddTask queues a function for asynchronous execution.
//   - RunWorkers starts the dispatcher if it is not already running.
//   - StopWorkers stops the current dispatcher loop.
//   - RestartWorkers stops the current dispatcher and starts a new one.
//
// Example:
//
//	workers.SetWorkersLimit(100)
//	workers.RunWorkers()
//	defer workers.StopWorkers()
//
//	workers.AddTask(func() {
//		// do work asynchronously
//	})
//
// AddTask blocks when the configured queue is full. RunWorkers is idempotent:
// calling it while the dispatcher is already running does not start another
// loop. StopWorkers stops task dispatching, but it does not cancel tasks that
// have already started.
package workers
