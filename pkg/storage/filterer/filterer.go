package filterer

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Filterer struct {
	storage storage.Storage
}

func NewFilterer(storage storage.Storage) *Filterer {
	return &Filterer{
		storage: storage,
	}
}

type filterFunc func(meta.Object) (Match, error)

// Find a single meta.Object of the given kind using the given filter
func (f *Filterer) Find(gvk schema.GroupVersionKind, filter BaseFilter) (meta.Object, error) {
	var results []Match
	var exactMatch Match

	// Fetch the sources, correct filtering method and if we're dealing with meta.APIType objects
	sources, filterFunc, metaObjects, err := f.parseFilter(gvk, filter)
	if err != nil {
		return nil, err
	}

	// Perform the filtering
	for _, object := range sources {
		if match, err := filterFunc(object); err != nil { // The filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			if match.Exact() {
				if exactMatch != nil {
					// We have multiple exact matches, the user has done something wrong
					return nil, filter.AmbiguousError([]Match{exactMatch, match})
				} else {
					exactMatch = match
				}
			} else {
				results = append(results, match)
			}
		}
	}

	var result meta.Object

	// If we have an exact result, select it
	if exactMatch != nil {
		result = exactMatch.Object()
	} else {
		if len(results) == 0 {
			return nil, filter.NonexistentError()
		} else if len(results) > 1 {
			return nil, filter.AmbiguousError(results)
		}

		result = results[0].Object()
	}

	// If we're filtering meta.APIType objects, load the full Object to be returned
	if metaObjects {
		return f.storage.Get(result.GroupVersionKind(), result.GetUID())
	}

	return result, nil
}

// Find all meta.Objects of the given kind using the given filter
func (f *Filterer) FindAll(gvk schema.GroupVersionKind, filter BaseFilter) ([]meta.Object, error) {
	var results []meta.Object

	// Fetch the sources, correct filtering method and if we're dealing with meta.APIType objects
	sources, filterFunc, metaObjects, err := f.parseFilter(gvk, filter)
	if err != nil {
		return nil, err
	}

	// Perform the filtering
	for _, object := range sources {
		if match, err := filterFunc(object); err != nil { // The filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			results = append(results, match.Object())
		}
	}

	// If we're filtering meta.APIType objects, load the full Objects to be returned
	if metaObjects {
		objects := make([]meta.Object, len(results))
		for i, result := range results {
			if objects[i], err = f.storage.Get(result.GroupVersionKind(), result.GetUID()); err != nil {
				return nil, err
			}
		}

		return objects, nil
	}

	return results, nil
}

func (f *Filterer) parseFilter(gvk schema.GroupVersionKind, filter BaseFilter) (sources []meta.Object, filterFunc filterFunc, metaObjects bool, err error) {
	// Parse ObjectFilters before MetaFilters, so ObjectFilters can embed MetaFilters
	if objectFilter, ok := filter.(ObjectFilter); ok {
		filterFunc = objectFilter.Filter
		sources, err = f.storage.List(gvk)
	} else if metaFilter, ok := filter.(MetaFilter); ok {
		filterFunc = metaFilter.FilterMeta
		sources, err = f.storage.ListMeta(gvk)
		metaObjects = true
	} else {
		err = fmt.Errorf("invalid filter type: %T", filter)
	}

	// Make sure the desired kind propagates down to the filter
	filter.SetKind(meta.Kind(gvk.Kind))

	return
}
