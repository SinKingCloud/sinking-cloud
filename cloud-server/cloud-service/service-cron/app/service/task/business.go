package task

// Start 启动任务
func (s *service) Start() {
	s.cron.Start()
}

// Close 关闭任务
func (s *service) Close() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}
