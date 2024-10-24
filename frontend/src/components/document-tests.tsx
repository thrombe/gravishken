import { FileText, NotepadText, Sheet, Presentation } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import * as server from '@common/server';
import * as types from '@common/types';
import { toast } from '@/hooks/use-toast';

const appMapping = {
  'docx': { icon: FileText, color: 'text-blue-600', appType: types.AppType.DOCX },
  'xlsx': { icon: Sheet, color: 'text-green-600', appType: types.AppType.XLSX },
  'pptx': { icon: Presentation, color: 'text-red-600', appType: types.AppType.PPTX },
  'txt': { icon: NotepadText, color: 'text-blue-600', appType: types.AppType.TXT },
};

interface DocumentTestsProps {
  testData: types.Test;
  handleFinishTest: () => void;
}

export default function DocumentTests({
  testData,
  handleFinishTest,
}: DocumentTestsProps) {
  const handleOpenApp = (app: types.Test) => {
    /// @ts-ignore
    let typ: types.AppType = appMapping[testData.Type].appType;
    server.server.send_message({
      Typ: types.Varient.OpenApp,
      Val: { Typ: typ, TestId: app.Id },
    });
  };

  const handleForceCloseApp = () => {
    server.server.send_message({
      Typ: types.Varient.QuitApp,
      Val: {},
    });
  };

  const handleSubmitWork = async () => {
    let resp = await fetch(server.base_url + "/get-user");
    let user: types.User = await resp.json()
    let submission: types.TestSubmission = {
      TestId: testData.Id,
      UserId: user.Id,
      TestInfo: {
        Type: testData.Type,
      },
    };
    const res = await fetch(server.base_url + "/submit-test", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(submission),
    })
    if (res.status === 400) {
      toast({
        variant: "destructive",
        title: "Failed to submit test",
        description: "Please ateast open the test first!",
      })
      return;
    }
    handleFinishTest()
  };

  // @ts-ignore
  const appConfig = appMapping[testData.Type];

  return (
    <Card className="h-full flex flex-col">
      <CardHeader>
        <CardTitle>
          {testData.TestName}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-grow">
        <div className="mb-4">
          <h3 className="text-lg font-semibold mb-2">Associated Apps:</h3>
          <div className="flex space-x-2 mb-4">

            <Button
              variant={"default"}
              onClick={() => handleOpenApp(testData)}
              className="flex items-center space-x-2"
            >
              <span>Open Associated App</span>
            </Button>

            <Button
              variant={"destructive"}
              onClick={() => handleForceCloseApp()}
              className="flex items-center space-x-2"
            >
              <span>Force Close App</span>
            </Button>

          </div>
        </div>
        <div className='flex-grow overflow-auto'>
          <img
            src={testData.FilePath}
            alt={`${testData.Id} Test`}
            className="w-full h-full object-contain"
          />
        </div>
        <Button onClick={handleSubmitWork} className="w-full mt-4">
          Submit Work
        </Button>
      </CardContent>
    </Card>
  );
};
