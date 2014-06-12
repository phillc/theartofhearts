require 'pty'

class Agent
  def initialize(options={})
    @command = options[:command]
    @compile = options[:compile]
  end

  def compile!
    return unless @compile
    command = "#{@compile}"
    puts "compiling: #{command}"
    puts `#{command}`
    raise "Compile error" unless $?.to_i == 0
  end

  def command
    %{ava play hearts --run="#{@command}"}
  end
end

desc "run ruby random agents"
task :agents, :number, :sleep do |t, args|
  args.with_defaults number: 4, sleep: 0.1

  agent = Agent.new(compile: 'make', command: 'bin/my_agent')

  agent.compile!

  pids = []
  args.number.to_i.times do |i|
    pids << fork do
      STDOUT.sync = true

      begin
        PTY.spawn(agent.command) do |stdin, stdout, pid|
          begin
            stdin.each { |line| puts "[#{i}] #{line}" }
          rescue Errno::EIO
            puts "Errno:EIO error, but this probably just means that the process has finished giving output"
          end
        end
      rescue PTY::ChildExited
        puts "The child process exited!"
      end
    end
    sleep args.sleep.to_f
  end

  pids.each{|pid| Process.wait(pid)}
end
