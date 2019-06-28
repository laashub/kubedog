package multitrack

import (
	"fmt"

	"github.com/flant/kubedog/pkg/tracker/daemonset"
	"github.com/flant/kubedog/pkg/tracker/replicaset"
	"k8s.io/client-go/kubernetes"
)

func (mt *multitracker) TrackDaemonSet(kube kubernetes.Interface, spec MultitrackSpec, opts MultitrackOptions) error {
	feed := daemonset.NewFeed()

	feed.OnAdded(func(ready bool) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetAdded(spec, feed, ready)
	})
	feed.OnReady(func() error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetReady(spec, feed)
	})
	feed.OnFailed(func(reason string) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetFailed(spec, feed, reason)
	})
	feed.OnEventMsg(func(msg string) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetEventMsg(spec, feed, msg)
	})
	feed.OnAddedReplicaSet(func(rs replicaset.ReplicaSet) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetAddedReplicaSet(spec, feed, rs)
	})
	feed.OnAddedPod(func(pod replicaset.ReplicaSetPod) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetAddedPod(spec, feed, pod)
	})
	feed.OnPodError(func(podError replicaset.ReplicaSetPodError) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetPodError(spec, feed, podError)
	})
	feed.OnPodLogChunk(func(chunk *replicaset.ReplicaSetPodLogChunk) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetPodLogChunk(spec, feed, chunk)
	})
	feed.OnStatusReport(func(status daemonset.DaemonSetStatus) error {
		mt.mux.Lock()
		defer mt.mux.Unlock()
		return mt.daemonsetStatusReport(spec, feed, status)
	})

	return feed.Track(spec.ResourceName, spec.Namespace, kube, opts.Options)
}

func (mt *multitracker) daemonsetAdded(spec MultitrackSpec, feed daemonset.Feed, ready bool) error {
	if ready {
		mt.DaemonSetsStatuses[spec.ResourceName] = feed.GetStatus()

		mt.displayResourceTrackerMessageF("ds", spec.ResourceName, "appears to be READY\n")

		return mt.handleResourceReadyCondition(mt.TrackingDaemonSets, spec)
	}

	mt.displayResourceTrackerMessageF("ds", spec.ResourceName, "added\n")

	return nil
}

func (mt *multitracker) daemonsetReady(spec MultitrackSpec, feed daemonset.Feed) error {
	mt.DaemonSetsStatuses[spec.ResourceName] = feed.GetStatus()

	mt.displayResourceTrackerMessageF("ds", spec.ResourceName, "become READY\n")

	return mt.handleResourceReadyCondition(mt.TrackingDaemonSets, spec)
}

func (mt *multitracker) daemonsetFailed(spec MultitrackSpec, feed daemonset.Feed, reason string) error {
	mt.displayResourceErrorF("ds", spec.ResourceName, "%s\n", reason)

	return mt.handleResourceFailure(mt.TrackingDaemonSets, "ds", spec, reason)
}

func (mt *multitracker) daemonsetEventMsg(spec MultitrackSpec, feed daemonset.Feed, msg string) error {
	mt.displayResourceEventF("ds", spec.ResourceName, "%s\n", msg)

	return nil
}

func (mt *multitracker) daemonsetAddedReplicaSet(spec MultitrackSpec, feed daemonset.Feed, rs replicaset.ReplicaSet) error {
	if !rs.IsNew {
		return nil
	}

	mt.displayResourceTrackerMessageF("ds", spec.ResourceName, "rs/%s added\n", rs.Name)

	return nil
}

func (mt *multitracker) daemonsetAddedPod(spec MultitrackSpec, feed daemonset.Feed, pod replicaset.ReplicaSetPod) error {
	if !pod.ReplicaSet.IsNew {
		return nil
	}

	mt.displayResourceTrackerMessageF("ds", spec.ResourceName, "po/%s added\n", pod.Name)

	return nil
}

func (mt *multitracker) daemonsetPodError(spec MultitrackSpec, feed daemonset.Feed, podError replicaset.ReplicaSetPodError) error {
	if !podError.ReplicaSet.IsNew {
		return nil
	}

	reason := fmt.Sprintf("po/%s container/%s: %s", podError.PodName, podError.ContainerName, podError.Message)

	mt.displayResourceErrorF("ds", spec.ResourceName, "%s\n", reason)

	return mt.handleResourceFailure(mt.TrackingDaemonSets, "ds", spec, reason)
}

func (mt *multitracker) daemonsetPodLogChunk(spec MultitrackSpec, feed daemonset.Feed, chunk *replicaset.ReplicaSetPodLogChunk) error {
	if !chunk.ReplicaSet.IsNew {
		return nil
	}

	mt.displayResourceLogChunk("ds", spec.ResourceName, podContainerLogChunkHeader(chunk.PodName, chunk.ContainerLogChunk), spec, chunk.ContainerLogChunk)

	return nil
}

func (mt *multitracker) daemonsetStatusReport(spec MultitrackSpec, feed daemonset.Feed, status daemonset.DaemonSetStatus) error {
	mt.DaemonSetsStatuses[spec.ResourceName] = status
	return nil
}
