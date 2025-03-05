package mock

import callback "dev.azure.com/2f-capital/go-packages/callback-client.git"

func Init() callback.Client {
	return &callbackClient{
		Service: Service{
			Status: callback.StatusActive,
			Events: make(map[string]*Event),
		},
	}
}
