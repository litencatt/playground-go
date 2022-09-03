class HelloInterface
 def hello
    raise NotImplementedError.new
 end
end

class TypeA < HelloInterface
  def hello
    puts "hello"
  end
end

class Foo
 def greeting(a)
   a.hello
 end
end

a = TypeA.new
f = Foo.new
f.greeting(a)
