package custom_helpers_service

import common_models "go_boilerplate_project/models/commons"

func (s *service) IsContextPresent(input ...common_models.IsContextPresentInput) bool {
	if len(input) == 0 || len(input) > 1 {
		return false
	}
	return input[0].Context != nil && input[0].CancelFunc != nil
}
