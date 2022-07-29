package StressClient

import "time"

type _SimClient struct {
	id  int
	cli _IClient
	cp  *_ClientPool

	isLogin    bool
	loginTime  time.Time
	logoutTime time.Time
}

func NewSimClient(id int, cp *_ClientPool) *_SimClient {
	sim := &_SimClient{
		id:         id,
		cp:         cp,
		cli:        cp.createClient(),
		isLogin:    false,
		logoutTime: time.Now(),
	}

	sim.cli.init(sim.cli, cp.loginAddr, cp.getAccount(id), "123123", []byte{})
	sim.cli.setDelegateToLoginSucceed(sim.onLoginSucceed)
	sim.cli.setDelegateToLoginFailed(sim.onLoginFailed)

	return sim
}

func (c *_SimClient) GetID() int {
	return c.id
}

func (c *_SimClient) CouldRecycle(leastTime time.Duration) bool {

	if !c.isLogin {
		return true
	}

	if time.Now().Sub(c.loginTime) < leastTime {
		return false
	}

	return true
}

func (c *_SimClient) CouldLogin() bool {
	if c.isLogin {
		return false
	}

	if time.Now().Sub(c.logoutTime) < 1*time.Second {
		return false
	}
	return true
}

func (c *_SimClient) Login() {
	c.cli.Login()
}

func (c *_SimClient) Logout() {
	c.isLogin = false
	c.cli.Logout()
	c.logoutTime = time.Now()
}

func (c *_SimClient) onLoginSucceed() {
	c.cp.onLoginSucceed(c)
	c.loginTime = time.Now()
	c.isLogin = true
}

func (c *_SimClient) onLoginFailed(err string) {
	c.cp.onLoginFailed(c)
}
